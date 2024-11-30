package br.com.dillmann.nginxignition.api.certificate.handler

import br.com.dillmann.nginxignition.api.certificate.CertificateConverter
import br.com.dillmann.nginxignition.core.certificate.command.GetCertificateCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.UuidAwareRequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond
import java.util.*

internal class GetCertificateByIdHandler(
    private val getCommand: GetCertificateCommand,
    private val converter: CertificateConverter,
): UuidAwareRequestHandler {
    override suspend fun handle(call: ApiCall, id: UUID) {
        val certificate = getCommand.getById(id)
        if (certificate != null) {
            val payload = converter.toResponse(certificate)
            call.respond(HttpStatus.OK, payload)
        } else {
            call.respond(HttpStatus.NOT_FOUND)
        }
    }
}
