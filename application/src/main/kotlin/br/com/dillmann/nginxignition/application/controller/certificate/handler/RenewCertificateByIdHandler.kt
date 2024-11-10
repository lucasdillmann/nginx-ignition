package br.com.dillmann.nginxignition.application.controller.certificate.handler

import br.com.dillmann.nginxignition.application.common.routing.template.IdAwareRequestHandler
import br.com.dillmann.nginxignition.application.controller.certificate.model.CertificateConverter
import br.com.dillmann.nginxignition.core.certificate.command.RenewCertificateCommand
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import java.util.UUID

class RenewCertificateByIdHandler(
    private val command: RenewCertificateCommand,
    private val converter: CertificateConverter,
): IdAwareRequestHandler {
    override suspend fun handle(call: RoutingCall, id: UUID) {
        val output = command.renewById(id)
        val responsePayload = converter.toResponse(output)
        val status = if (output.success) HttpStatusCode.OK else HttpStatusCode.BadRequest
        call.respond(status, responsePayload)
    }
}
