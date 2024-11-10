package br.com.dillmann.nginxignition.application.controller.certificate.handler

import br.com.dillmann.nginxignition.application.common.routing.RequestHandler
import br.com.dillmann.nginxignition.application.controller.certificate.model.CertificateConverter
import br.com.dillmann.nginxignition.core.certificate.command.GetCertificateCommand
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import java.util.*

class GetCertificateByIdHandler(
    private val getCommand: GetCertificateCommand,
    private val converter: CertificateConverter,
): RequestHandler {
    override suspend fun handle(call: RoutingCall) {
        val certificateId = runCatching { call.request.pathVariables["id"].let(UUID::fromString) }.getOrNull()
        if (certificateId == null) {
            call.respond(HttpStatusCode.BadRequest)
            return
        }

        val certificate = getCommand.getById(certificateId)
        if (certificate != null) {
            val payload = converter.toResponse(certificate)
            call.respond(HttpStatusCode.OK, payload)
        } else {
            call.respond(HttpStatusCode.NotFound)
        }
    }
}
