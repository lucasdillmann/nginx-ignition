package br.com.dillmann.nginxignition.application

import br.com.dillmann.nginxignition.application.common.configuration.configureHttp
import br.com.dillmann.nginxignition.application.common.configuration.configureKoin
import br.com.dillmann.nginxignition.application.common.configuration.configureLifecycle
import br.com.dillmann.nginxignition.application.common.configuration.configureRoutes
import br.com.dillmann.nginxignition.core.common.log.logger
import io.ktor.server.application.*
import io.ktor.server.netty.*
import kotlinx.coroutines.runBlocking

fun main(args: Array<String>) {
    EngineMain.main(args)
}

fun Application.module() {
    logger("Application").info("Welcome to nginx ignition")

    runBlocking {
        configureKoin()
        configureHttp()
        configureRoutes()
        configureLifecycle()
    }
}