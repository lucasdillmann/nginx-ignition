package br.com.dillmann.nginxsidewheel.core.nginx.configuration

import br.com.dillmann.nginxsidewheel.core.host.Host

internal interface NginxConfigurationFileProvider {
    data class Output(
        val name: String,
        val contents: String,
    )

    suspend fun provide(basePath: String, hosts: List<Host>): List<Output>
}
