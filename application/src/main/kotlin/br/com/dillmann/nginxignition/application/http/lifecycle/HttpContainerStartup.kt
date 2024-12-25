package br.com.dillmann.nginxignition.application.http.lifecycle

import br.com.dillmann.nginxignition.application.http.HttpContainer
import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand

internal class HttpContainerStartup(private val container: HttpContainer): StartupCommand {
    override suspend fun execute() {
        container.start()
    }
}
