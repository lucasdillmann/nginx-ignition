package br.com.dillmann.nginxignition.application.controller.certificate

import br.com.dillmann.nginxignition.application.common.rbac.requireAuthentication
import br.com.dillmann.nginxignition.application.common.routing.*
import br.com.dillmann.nginxignition.application.controller.certificate.handler.*
import io.ktor.server.application.*
import io.ktor.server.routing.*
import org.koin.ktor.ext.inject

fun Application.certificateRoutes() {
    val listProvidersHandler by inject<ListAvailableProvidersHandler>()
    val issueHandler by inject<IssueCertificateHandler>()
    val listHandler by inject<ListCertificatesHandler>()
    val deleteByIdHandler by inject<DeleteCertificateByIdHandler>()
    val getByIdHandler by inject<GetCertificateByIdHandler>()
    val renewByIdHandler by inject<RenewCertificateByIdHandler>()

    routing {
        requireAuthentication {
            route("/api/certificates") {
                get(listHandler)
                get("/{id}", getByIdHandler)
                delete("/{id}", deleteByIdHandler)
                post("/{id}/renew", renewByIdHandler)
                post("/issue", issueHandler)
                get("/available-providers", listProvidersHandler)
            }
        }
    }
}
