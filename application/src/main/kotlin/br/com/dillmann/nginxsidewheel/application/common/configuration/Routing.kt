package br.com.dillmann.nginxsidewheel.application.common.configuration

import br.com.dillmann.nginxsidewheel.application.controller.frontend.frontendRoutes
import br.com.dillmann.nginxsidewheel.application.controller.host.hostRoutes
import br.com.dillmann.nginxsidewheel.application.controller.nginx.nginxRoutes
import io.ktor.server.application.*
import io.ktor.server.routing.*

fun Application.configureRoutes() {
    routing {
        hostRoutes()
        nginxRoutes()
        frontendRoutes()
    }
}
