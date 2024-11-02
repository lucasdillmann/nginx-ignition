package br.com.dillmann.nginxsidewheel.application

import br.com.dillmann.nginxsidewheel.application.common.configuration.configureHttp
import br.com.dillmann.nginxsidewheel.application.common.configuration.configureKoin
import br.com.dillmann.nginxsidewheel.application.common.configuration.configureRoutes
import br.com.dillmann.nginxsidewheel.application.common.configuration.runStartupCommands
import io.ktor.server.application.*
import io.ktor.server.netty.*

fun main(args: Array<String>) {
    EngineMain.main(args)
}

fun Application.module() {
    configureKoin()
    configureHttp()
    configureRoutes()
    runStartupCommands()
}
