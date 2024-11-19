package br.com.dillmann.nginxignition.application.configuration

import br.com.dillmann.nginxignition.application.api.apiRoutes
import br.com.dillmann.nginxignition.application.frontend.frontendRoutes
import io.ktor.server.application.*
import io.ktor.server.routing.*

fun Application.configureRoutes() {
    routing {
        apiRoutes()
        frontendRoutes()
    }
}
