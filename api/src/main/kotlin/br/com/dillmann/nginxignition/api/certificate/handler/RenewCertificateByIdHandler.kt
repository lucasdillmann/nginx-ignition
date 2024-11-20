package br.com.dillmann.nginxignition.api.certificate.handler

import br.com.dillmann.nginxignition.api.certificate.model.CertificateConverter
import br.com.dillmann.nginxignition.core.certificate.command.RenewCertificateCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.IdAwareRequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond
import java.util.UUID

internal class RenewCertificateByIdHandler(
    private val command: RenewCertificateCommand,
    private val converter: CertificateConverter,
): IdAwareRequestHandler {
    override suspend fun handle(call: ApiCall, id: UUID) {
        val output = command.renewById(id)
        val responsePayload = converter.toResponse(output)
        val status = if (output.success) HttpStatus.OK else HttpStatus.BAD_REQUEST
        call.respond(status, responsePayload)
    }
}
