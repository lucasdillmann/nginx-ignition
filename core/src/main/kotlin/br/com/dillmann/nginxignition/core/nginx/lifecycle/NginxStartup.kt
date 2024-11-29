package br.com.dillmann.nginxignition.core.nginx.lifecycle

import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxignition.core.nginx.NginxService

internal class NginxStartup(private val service: NginxService): StartupCommand {
    @Suppress("MagicNumber")
    override val priority = 500
    private val logger = logger<NginxStartup>()

    override suspend fun execute() {
        try {
            service.start()
        } catch (ex: Exception) {
            logger.warn("Nginx failed to start automatically. Please check your hosts configurations.", ex)
        }
    }
}
