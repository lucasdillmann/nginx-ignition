package br.com.dillmann.nginxsidewheel.application.configuration

import br.com.dillmann.nginxsidewheel.application.controller.host.hostRoutes
import io.ktor.server.application.*
import io.ktor.server.routing.*

fun Application.configureRoutes() {
    routing {
        hostRoutes()
    }
}
