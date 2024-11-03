package br.com.dillmann.nginxignition.core.nginx.lifecycle

import br.com.dillmann.nginxignition.core.common.lifecycle.ShutdownCommand
import br.com.dillmann.nginxignition.core.nginx.NginxService

internal class NginxShutdown(private val service: NginxService): ShutdownCommand {
    override val priority = 500

    override suspend fun execute() {
        if (!service.isRunning()) return
        service.stop()
    }
}
