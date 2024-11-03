package br.com.dillmann.nginxsidewheel.core.nginx.lifecycle

import br.com.dillmann.nginxsidewheel.core.common.log.logger
import br.com.dillmann.nginxsidewheel.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxsidewheel.core.nginx.NginxService

internal class NginxStartup(private val service: NginxService): StartupCommand {
    private val logger = logger<NginxStartup>()
    override val priority = 500

    override suspend fun execute() {
        try {
            service.start()
        } catch (ex: Exception) {
            logger.warn("Nginx failed to start automatically. Please check your hosts configurations.", ex)
        }
    }
}
