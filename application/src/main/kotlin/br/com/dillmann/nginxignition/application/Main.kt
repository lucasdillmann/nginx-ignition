package br.com.dillmann.nginxignition.application

import kotlinx.coroutines.runBlocking
import org.slf4j.LoggerFactory
import kotlin.system.exitProcess

private val LOGGER = LoggerFactory.getLogger("Main")

fun main() {
    Thread.currentThread().setUncaughtExceptionHandler { _, ex ->
        LOGGER.error("Application startup failed", ex)
        exitProcess(1)
    }

    runBlocking {
        Application().boot()
    }
}
