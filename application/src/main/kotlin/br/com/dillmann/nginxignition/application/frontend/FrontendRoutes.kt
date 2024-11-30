package br.com.dillmann.nginxignition.application.frontend

import io.ktor.server.application.*
import io.ktor.server.http.content.*
import io.ktor.server.routing.*

fun Application.frontendRoutes() {
    routing {
        staticResources(
            remotePath = "/",
            basePackage = "nginx-ignition/frontend",
            index = "index.html",
        ) {
            default("index.html")
        }
    }
}
