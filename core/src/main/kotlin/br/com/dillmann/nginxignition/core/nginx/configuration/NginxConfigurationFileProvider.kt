package br.com.dillmann.nginxignition.core.nginx.configuration

import br.com.dillmann.nginxignition.core.host.Host

internal fun interface NginxConfigurationFileProvider {
    data class Output(
        val name: String,
        val contents: String,
    )

    suspend fun provide(basePath: String, hosts: List<Host>): List<Output>
}
