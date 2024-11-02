package br.com.dillmann.nginxsidewheel.application.common.configuration

import br.com.dillmann.nginxsidewheel.core.common.startup.StartupCommand
import io.ktor.server.application.*
import org.koin.ktor.ext.getKoin

suspend fun Application.runStartupCommands() {
    getKoin()
        .getAll<StartupCommand>()
        .sortedBy { it.priority }
        .forEach { it.execute() }
}
