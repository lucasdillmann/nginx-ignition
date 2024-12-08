package br.com.dillmann.nginxignition.application

import br.com.dillmann.nginxignition.application.configuration.*
import br.com.dillmann.nginxignition.core.common.log.logger
import io.ktor.server.application.*
import io.ktor.server.netty.*
import kotlinx.coroutines.runBlocking

fun main(args: Array<String>) {
    EngineMain.main(args)
}

fun Application.module() {
    val logger = logger("Application")
    logger.info("Welcome to nginx ignition")

    try {
        runBlocking {
            configureKoin()
            configureRbac()
            configureHttp()
            configureRoutes()
            configureExceptionHandling()
            configureLifecycle()
        }
    } catch (ex: Exception) {
        logger.error("Application startup failed", ex)
        Runtime.getRuntime().halt(1)
    }
}
