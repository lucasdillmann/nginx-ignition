package br.com.dillmann.nginxignition.application.controller.certificate

import br.com.dillmann.nginxignition.application.common.rbac.requireAuthentication
import br.com.dillmann.nginxignition.application.common.routing.*
import br.com.dillmann.nginxignition.application.controller.certificate.handler.IssueCertificateHandler
import br.com.dillmann.nginxignition.application.controller.certificate.handler.ListAvailableProvidersHandler
import io.ktor.server.application.*
import io.ktor.server.routing.*
import org.koin.ktor.ext.inject

fun Application.certificateRoutes() {
    val listProvidersHandler by inject<ListAvailableProvidersHandler>()
    val issueHandler by inject<IssueCertificateHandler>()

    routing {
        requireAuthentication {
            route("/api/certificates") {
                get { TODO() }
                delete("/{id}") { TODO() }
                post("/{id}/renew") { TODO() }
                post("/issue", issueHandler)
                get("/available-providers", listProvidersHandler)
            }
        }
    }
}
