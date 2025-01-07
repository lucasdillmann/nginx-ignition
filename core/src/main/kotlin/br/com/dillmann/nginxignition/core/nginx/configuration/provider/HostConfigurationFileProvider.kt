package br.com.dillmann.nginxignition.core.nginx.configuration.provider

import br.com.dillmann.nginxignition.core.host.Host
import br.com.dillmann.nginxignition.core.integration.IntegrationService
import br.com.dillmann.nginxignition.core.nginx.configuration.NginxConfigurationFileProvider
import br.com.dillmann.nginxignition.core.settings.SettingsRepository
import java.net.URI

internal class HostConfigurationFileProvider(
    private val integrationService: IntegrationService,
    private val settingsService: SettingsRepository,
): NginxConfigurationFileProvider {
    override suspend fun provide(basePath: String, hosts: List<Host>): List<NginxConfigurationFileProvider.Output> =
        hosts.filter { it.enabled }.map { buildHost(basePath, it) }

    private suspend fun buildHost(basePath: String, host: Host): NginxConfigurationFileProvider.Output {
        val routes = host
            .routes
            .filter { it.enabled }
            .sortedBy { it.priority }
            .map { buildRoute(host, it, basePath) }
            .joinToString("\n")
        val serverNames =
            if (host.defaultServer) "server_name _;"
            else host.domainNames?.joinToString("\n") { "server_name $it;" } ?: ""
        val httpsRedirect =
            if (host.featureSet.redirectHttpToHttps)
                """
                    if (${'$'}scheme = "http") {
                        return 301 https://${'$'}server_name${'$'}request_uri;
                    }
                """.trimIndent()
            else ""
        val http2 =
            if (host.featureSet.http2Support) "http2 on;"
            else ""

        val bindings =
            if (host.useGlobalBindings) settingsService.get().globalBindings
            else host.bindings
        val contents = bindings
            .map { buildBinding(basePath, host, it, routes, serverNames, httpsRedirect, http2) }
            .joinToString("\n")

        return NginxConfigurationFileProvider.Output(
            name = "host-${host.id}.conf",
            contents = contents,
        )
    }

    private suspend fun buildBinding(
        basePath: String,
        host: Host,
        binding: Host.Binding,
        routes: String,
        serverNames: String,
        httpsRedirect: String,
        http2: String,
    ): String {
        val listen =
            when(binding.type) {
                Host.BindingType.HTTP ->
                    "listen ${binding.ip}:${binding.port} ${buildBindingAdditionalParams(host)};"
                Host.BindingType.HTTPS ->
                    """
                        listen ${binding.ip}:${binding.port} ssl ${buildBindingAdditionalParams(host)};
                        ssl_certificate $basePath/config/certificate-${binding.certificateId}.pem;
                        ssl_certificate_key $basePath/config/certificate-${binding.certificateId}.pem;
                        ssl_protocols TLSv1.2 TLSv1.3;
                        ssl_ciphers HIGH:!aNULL:!MD5;
                    """.trimIndent()
            }

        val conditionalHttpsRedirect =
            if (binding.type == Host.BindingType.HTTP) httpsRedirect
            else ""

        val settings = settingsService.get().nginx
        val logs = settings.logs

        return """
            server {
                root /dev/null;
                access_log ${
                    if(logs.accessLogsEnabled) "$basePath/logs/host-${host.id}.access.log" 
                    else "off"
                };
                error_log ${
                    if(logs.errorLogsEnabled) 
                        "$basePath/logs/host-${host.id}.error.log ${logs.errorLogsLevel.name.lowercase()}" 
                    else 
                        "off"
                };
                gzip ${flag(settings.gzipEnabled)};
                client_max_body_size ${settings.maximumBodySizeMb}M;
                ${flag(
                    host.accessListId != null, 
                    "include $basePath/config/access-list-${host.accessListId}.conf;", 
                    ""
                )}
                
                $conditionalHttpsRedirect
                $http2
                $listen
                $serverNames
                $routes
            }
        """.trimIndent()
    }

    private fun buildBindingAdditionalParams(host: Host): String {
        var additionalParams = ""
        if (host.defaultServer)
            additionalParams += " default_server"

        return additionalParams
    }

    private suspend fun buildRoute(
        host: Host,
        route: Host.Route,
        basePath: String,
    ): String =
        when (route.type) {
            Host.RouteType.STATIC_RESPONSE -> buildStaticResponseRoute(route, host.featureSet, basePath)
            Host.RouteType.PROXY -> buildProxyRoute(route, host.featureSet, basePath)
            Host.RouteType.REDIRECT -> buildRedirectRoute(route, host.featureSet, basePath)
            Host.RouteType.INTEGRATION -> buildIntegrationRoute(route, host.featureSet, basePath)
            Host.RouteType.SOURCE_CODE -> buildSourceCodeRoute(host, route, basePath)
        }

    private fun buildSourceCodeRoute(host: Host, route: Host.Route, basePath: String): String {
        val (_, priority, _, _, sourcePath) = route
        val (language, code, mainFunction) = route.sourceCode!!
        val (headerBlock, routeBlock) = when (language) {
            Host.SourceCodeLanguage.JAVASCRIPT ->
                Pair(
                    "js_import route_$priority from $basePath/config/host-${host.id}-route-$priority.js;",
                    "js_content route_$priority.$mainFunction;",
                )
            Host.SourceCodeLanguage.LUA ->
                Pair(
                    "",
                    """
                        content_by_lua_block {
                            $code
                        }
                    """.trimIndent(),
                )
        }

        return """
            $headerBlock
            location $sourcePath {
                $routeBlock
                ${buildRouteFeatures(host.featureSet)}
                ${buildRouteSettings(route, basePath)}
            }
        """.trimIndent()
    }

    private fun buildStaticResponseRoute(route: Host.Route, features: Host.FeatureSet, basePath: String): String {
        val (statusCode, payload, headers) = route.response!!
        val headerInstructions = headers
            .entries
            .joinToString(separator = "\n") { (key, value) -> "add_header \"${key.scape()}\" \"${value.scape()}\";" }
        val escapedPayload = (payload ?: "").scape()

        return """
            location ${route.sourcePath} {
                $headerInstructions
                return $statusCode "$escapedPayload";
                ${buildRouteFeatures(features)}
                ${buildRouteSettings(route, basePath)}
            }
        """.trimIndent()
    }

    private fun buildProxyRoute(route: Host.Route, features: Host.FeatureSet, basePath: String) =
        """
           location ${route.sourcePath} {
                ${buildProxyPass(route)}
                ${buildRouteFeatures(features)}
                ${buildRouteSettings(route, basePath)}
           } 
        """.trimIndent()

    private suspend fun buildIntegrationRoute(route: Host.Route, features: Host.FeatureSet, basePath: String): String {
        val (integrationId, integrationOption) = route.integration!!
        val proxyUrl = integrationService.getIntegrationOptionUrl(integrationId, integrationOption)
        return """
           location ${route.sourcePath} {
                ${buildProxyPass(route, proxyUrl)}
                ${buildRouteFeatures(features)}
                ${buildRouteSettings(route, basePath)}
           }
        """.trimIndent()
    }

    private fun buildRedirectRoute(route: Host.Route, features: Host.FeatureSet, basePath: String) =
        """
           location ${route.sourcePath} {
                return ${route.redirectCode!!} ${route.targetUri!!};
                ${buildRouteFeatures(features)}
                ${buildRouteSettings(route, basePath)}
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

    private fun buildProxyPass(route: Host.Route, uri: String = route.targetUri!!): String {
        val builder = StringBuilder().appendLine("proxy_pass $uri;")

        if (route.settings.keepOriginalDomainName)
            builder.appendLine("proxy_set_header Host ${URI(uri).host};")

        return builder.toString()
    }

    private fun buildRouteSettings(route: Host.Route, basePath: String): String {
        val builder = StringBuilder()
        with (route.settings) {
            if (proxySslServerName)
                builder.appendLine("proxy_ssl_server_name on;")

            if (includeForwardHeaders) {
                builder.appendLine("proxy_set_header x-forwarded-for \$proxy_add_x_forwarded_for;")
                builder.appendLine("proxy_set_header x-forwarded-host \$host;")
                builder.appendLine("proxy_set_header x-forwarded-proto \$scheme;")
                builder.appendLine("proxy_set_header x-forwarded-scheme \$scheme;")
                builder.appendLine("proxy_set_header x-forwarded-port \$server_port;")
                builder.appendLine("proxy_set_header x-real-ip \$remote_addr;")
            }

            if (!custom.isNullOrBlank())
                builder.appendLine(custom)
        }

        if (route.accessListId != null)
            builder.appendLine("include $basePath/config/access-list-${route.accessListId}.conf;")

        return builder.toString()
    }

    private fun flag(
        enabled: Boolean,
        trueValue: String = "on",
        falseValue: String = "off",
    ) = if (enabled) trueValue else falseValue
}
