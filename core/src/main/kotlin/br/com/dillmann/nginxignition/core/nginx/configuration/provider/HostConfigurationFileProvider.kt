package br.com.dillmann.nginxignition.core.nginx.configuration.provider

import br.com.dillmann.nginxignition.core.host.Host
import br.com.dillmann.nginxignition.core.nginx.configuration.NginxConfigurationFileProvider

internal class HostConfigurationFileProvider: NginxConfigurationFileProvider {
    override suspend fun provide(basePath: String, hosts: List<Host>): List<NginxConfigurationFileProvider.Output> =
        hosts.filter { it.enabled }.map { buildHost(basePath, it) }

    private fun buildHost(basePath: String, host: Host): NginxConfigurationFileProvider.Output {
        val routes = host.routes.sortedBy { it.priority }.joinToString(separator = "\n") { buildRoute(it) }
        val bindings = host.bindings.joinToString(separator = "\n") { buildBinding(it, host.default) }
        val serverNames =
            if (host.default) "server_name _;"
            else host.domainNames.joinToString(separator = "\n") { "server_name $it;" }

        val contents = """
            server {
                root /dev/null;
                access_log $basePath/logs/host-${host.id}.access.log;
                error_log $basePath/logs/host-${host.id}.error.log;
                gzip on;
                client_max_body_size 1024G;
                
                $bindings
                $serverNames
                $routes
            }
        """.trimIndent()

        // TODO: Implement host features (such as websocket support)
        return NginxConfigurationFileProvider.Output(
            name = "host-${host.id}.conf",
            contents = contents,
        )
    }

    private fun buildBinding(binding: Host.Binding, default: Boolean): String =
        when (binding.type) {
            Host.BindingType.HTTP -> buildHttpBinding(binding, default)
            Host.BindingType.HTTPS -> buildHttpsBinding(binding, default)
        }

    private fun buildHttpBinding(binding: Host.Binding, default: Boolean): String =
        "listen ${binding.ip}:${binding.port}${if (default) " default_server" else ""};"

    private fun buildHttpsBinding(binding: Host.Binding, default: Boolean): String =
        TODO()

    private fun buildRoute(route: Host.Route): String =
        when (route.type) {
            Host.RouteType.STATIC_RESPONSE -> buildStaticResponseRoute(route)
            Host.RouteType.PROXY -> buildProxyRoute(route)
            Host.RouteType.REDIRECT -> buildRedirectRoute(route)
        }

    private fun buildStaticResponseRoute(route: Host.Route): String {
        val (statusCode, payload, headers) = route.response!!
        val headerInstructions = headers
            .flatMap { (key, values) -> values.map { key to it } }
            .joinToString(separator = "\n") { (key, value) -> "add_header \"${key.scape()}\" \"${value.scape()}\";" }
        val escapedPayload = (payload ?: "").scape()

        return """
            location ${route.sourcePath} {
                $headerInstructions
                return $statusCode "$escapedPayload";
            }
        """.trimIndent()
    }

    private fun buildProxyRoute(route: Host.Route) =
        """
           location ${route.sourcePath} {
                proxy_pass ${route.targetUri!!};
                proxy_set_header Host $\http_host;
                proxy_set_header X-Real-IP $\remote_addr;
                proxy_set_header X-Forwarded-For $\proxy_add_x_forwarded_for;
                proxy_set_header X-Forwarded-Proto $\scheme;
           } 
        """.trimIndent()

    private fun buildRedirectRoute(route: Host.Route) =
        """
           location ${route.sourcePath} {
                return ${route.redirectCode!!} ${route.targetUri!!};
           } 
        """.trimIndent()

    private fun String.scape() =
        replace("\"", "\\\"")
}
