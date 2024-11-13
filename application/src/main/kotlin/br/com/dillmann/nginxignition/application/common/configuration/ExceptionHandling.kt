package br.com.dillmann.nginxignition.application.common.configuration

import br.com.dillmann.nginxignition.application.common.exception.ConsistencyExceptionHandler
import br.com.dillmann.nginxignition.core.common.validation.ConsistencyException
import io.ktor.server.application.*
import io.ktor.server.plugins.statuspages.*
import org.koin.ktor.ext.inject

fun Application.configureExceptionHandling() {
    val consistencyExceptionHandler by inject<ConsistencyExceptionHandler>()

    install(StatusPages) {
        exception<ConsistencyException>(consistencyExceptionHandler::handle)
    }
}
