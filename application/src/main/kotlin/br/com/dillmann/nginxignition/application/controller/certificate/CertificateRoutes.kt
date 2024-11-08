package br.com.dillmann.nginxignition.application.controller.certificate

import br.com.dillmann.nginxignition.application.common.rbac.requireAuthentication
import br.com.dillmann.nginxignition.application.controller.certificate.handler.ListAvailableProvidersHandler
import io.ktor.server.application.*
import io.ktor.server.routing.*
import org.koin.ktor.ext.inject

fun Application.certificateRoutes() {
    val listProvidersHandler by inject<ListAvailableProvidersHandler>()

    routing {
        requireAuthentication {
            route("/api/certificates") {
                get { TODO() }
                delete("/{id}") { TODO() }
                post("/{id}/renew") { TODO() }
                post("/issue") { TODO() }
                get("/available-providers") { listProvidersHandler.handle(call) }
            }
        }
    }
}
