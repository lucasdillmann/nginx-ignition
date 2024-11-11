package br.com.dillmann.nginxignition.core.nginx.configuration.provider

import br.com.dillmann.nginxignition.core.host.Host
import br.com.dillmann.nginxignition.core.nginx.configuration.NginxConfigurationFileProvider

internal class HostConfigurationFileProvider: NginxConfigurationFileProvider {
    override suspend fun provide(basePath: String, hosts: List<Host>): List<NginxConfigurationFileProvider.Output> =
        hosts.filter { it.enabled }.map { buildHost(basePath, it) }

    private fun buildHost(basePath: String, host: Host): NginxConfigurationFileProvider.Output {
        val routes = host.routes.sortedBy { it.priority }.joinToString("\n") { buildRoute(it, host.featureSet) }
        val bindings = host.bindings.joinToString("\n") { buildBinding(it, host.default, basePath, host.featureSet) }
        val serverNames =
            if (host.default) "server_name _;"
            else host.domainNames.joinToString("\n") { "server_name $it;" }
        val httpsRedirect =
            if (host.featureSet.redirectHttpToHttps)
                """
                    if (${'$'}scheme = "http") {
                        return 301 https://${'$'}server_name${'$'}request_uri;
                    }
                """.trimIndent()
            else ""

        val contents = """
            server {
                root /dev/null;
                access_log $basePath/logs/host-${host.id}.access.log;
                error_log $basePath/logs/host-${host.id}.error.log;
                gzip on;
                client_max_body_size 1024G;
                
                $httpsRedirect
                $bindings
                $serverNames
                $routes
            }
        """.trimIndent()

        return NginxConfigurationFileProvider.Output(
            name = "host-${host.id}.conf",
            contents = contents,
        )
    }

    private fun buildBinding(
        binding: Host.Binding,
        default: Boolean,
        basePath: String,
        features: Host.FeatureSet,
    ): String =
        when (binding.type) {
            Host.BindingType.HTTP -> buildHttpBinding(binding, default, features)
            Host.BindingType.HTTPS -> buildHttpsBinding(binding, default, basePath, features)
        }

    private fun buildHttpBinding(
        binding: Host.Binding,
        default: Boolean,
        features: Host.FeatureSet,
    ): String =
        "listen ${binding.ip}:${binding.port} ${buildBindingAdditionalParams(default, features)};"

    private fun buildHttpsBinding(
        binding: Host.Binding,
        default: Boolean,
        basePath: String,
        features: Host.FeatureSet,
    ): String =
        """
            listen ${binding.ip}:${binding.port} ssl ${buildBindingAdditionalParams(default, features)};
            ssl_certificate $basePath/config/certificate-${binding.certificateId}.pem;
            ssl_certificate_key $basePath/config/certificate-${binding.certificateId}.pem;
            ssl_protocols TLSv1.2 TLSv1.3;
            ssl_ciphers HIGH:!aNULL:!MD5;
        """.trimIndent()

    private fun buildBindingAdditionalParams(default: Boolean, features: Host.FeatureSet): String {
        var additionalParams = ""
        if (features.http2Support)
            additionalParams += " http2"
        if (default)
            additionalParams += " default_server"

        return additionalParams
    }


    private fun buildRoute(route: Host.Route, features: Host.FeatureSet): String =
        when (route.type) {
            Host.RouteType.STATIC_RESPONSE -> buildStaticResponseRoute(route, features)
            Host.RouteType.PROXY -> buildProxyRoute(route, features)
            Host.RouteType.REDIRECT -> buildRedirectRoute(route, features)
        }

    private fun buildStaticResponseRoute(route: Host.Route, features: Host.FeatureSet): String {
        val (statusCode, payload, headers) = route.response!!
        val headerInstructions = headers
            .flatMap { (key, values) -> values.map { key to it } }
            .joinToString(separator = "\n") { (key, value) -> "add_header \"${key.scape()}\" \"${value.scape()}\";" }
        val escapedPayload = (payload ?: "").scape()

        return """
            location ${route.sourcePath} {
                $headerInstructions
                return $statusCode "$escapedPayload";
                ${buildRouteFeatures(features)}
            }
        """.trimIndent()
    }

    private fun buildProxyRoute(route: Host.Route, features: Host.FeatureSet) =
        """
           location ${route.sourcePath} {
                proxy_pass ${route.targetUri!!};
                proxy_set_header Host ${'$'}http_host;
                proxy_set_header X-Real-IP ${'$'}remote_addr;
                proxy_set_header X-Forwarded-For ${'$'}proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto ${'$'}scheme;
                ${buildRouteFeatures(features)}
           } 
        """.trimIndent()

    private fun buildRedirectRoute(route: Host.Route, features: Host.FeatureSet) =
        """
           location ${route.sourcePath} {
                return ${route.redirectCode!!} ${route.targetUri!!};
                ${buildRouteFeatures(features)}
           } 
        """.trimIndent()

    private fun buildRouteFeatures(features: Host.FeatureSet): String =
        if (features.websocketsSupport)
            """
                proxy_http_version 1.1;
                proxy_set_header Upgrade ${'$'}http_upgrade;
                proxy_set_header Connection "upgrade";
            """.trimIndent()
        else ""

    private fun String.scape() =
        replace("\"", "\\\"")
}
