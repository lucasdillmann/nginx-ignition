package br.com.dillmann.nginxsidewheel.core.nginx

import br.com.dillmann.nginxsidewheel.core.common.log.logger
import br.com.dillmann.nginxsidewheel.core.common.provider.ConfigurationProvider
import br.com.dillmann.nginxsidewheel.core.host.HostService

internal class NginxConfigurationFiles(
    private val hostService: HostService,
    private val configurationProvider: ConfigurationProvider,
) {
    private val logger = logger<NginxConfigurationFiles>()

    suspend fun generate() {
        val hosts = hostService.getAll()
        logger.info("Rebuilding nginx configuration files for ${hosts.size} hosts")

        // TODO
    }
}
