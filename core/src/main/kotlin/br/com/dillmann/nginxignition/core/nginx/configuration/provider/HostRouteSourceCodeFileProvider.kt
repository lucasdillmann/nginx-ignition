package br.com.dillmann.nginxignition.core.nginx.configuration.provider

import br.com.dillmann.nginxignition.core.host.Host
import br.com.dillmann.nginxignition.core.nginx.configuration.NginxConfigurationFileProvider

internal class HostRouteSourceCodeFileProvider: NginxConfigurationFileProvider {
    override suspend fun provide(basePath: String, hosts: List<Host>): List<NginxConfigurationFileProvider.Output> =
        hosts.filter { it.enabled }.flatMap { buildSourceCodeFiles(it) }

    private fun buildSourceCodeFiles(host: Host): List<NginxConfigurationFileProvider.Output> =
        host.routes
            .filter {
                it.enabled &&
                    it.type == Host.RouteType.SOURCE_CODE &&
                    it.sourceCode?.language == Host.SourceCodeLanguage.JAVASCRIPT
            }
            .map {
                NginxConfigurationFileProvider.Output(
                    name = "host-${host.id}-route-${it.priority}.js",
                    contents = it.sourceCode!!.code,
                )
            }
}
