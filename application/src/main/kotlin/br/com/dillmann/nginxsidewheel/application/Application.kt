package br.com.dillmann.nginxsidewheel.application

import br.com.dillmann.nginxsidewheel.application.common.configuration.configureHttp
import br.com.dillmann.nginxsidewheel.application.common.configuration.configureKoin
import br.com.dillmann.nginxsidewheel.application.common.configuration.configureRoutes
import br.com.dillmann.nginxsidewheel.application.common.configuration.runStartupCommands
import br.com.dillmann.nginxsidewheel.core.common.log.logger
import io.ktor.server.application.*
import io.ktor.server.netty.*
import kotlinx.coroutines.runBlocking

fun main(args: Array<String>) {
    EngineMain.main(args)
}

fun Application.module() {
    logger("Application").info("Welcome to nginx sidewheel")

    runBlocking {
        configureKoin()
        configureHttp()
        configureRoutes()
        runStartupCommands()
    }
}
