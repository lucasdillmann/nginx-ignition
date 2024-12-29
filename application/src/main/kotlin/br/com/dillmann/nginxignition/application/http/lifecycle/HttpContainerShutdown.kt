package br.com.dillmann.nginxignition.application.http.lifecycle

import br.com.dillmann.nginxignition.application.http.HttpContainer
import br.com.dillmann.nginxignition.core.common.lifecycle.ShutdownCommand

internal class HttpContainerShutdown(private val container: HttpContainer): ShutdownCommand {
    override suspend fun execute() {
        container.stop()
    }
}
