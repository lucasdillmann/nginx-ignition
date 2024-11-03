package br.com.dillmann.nginxignition.application.common.configuration

import br.com.dillmann.nginxignition.application.controller.frontend.frontendRoutes
import br.com.dillmann.nginxignition.application.controller.host.hostRoutes
import br.com.dillmann.nginxignition.application.controller.nginx.nginxRoutes
import io.ktor.server.application.*
import io.ktor.server.routing.*

fun Application.configureRoutes() {
    routing {
        hostRoutes()
        nginxRoutes()
        frontendRoutes()
    }
}
