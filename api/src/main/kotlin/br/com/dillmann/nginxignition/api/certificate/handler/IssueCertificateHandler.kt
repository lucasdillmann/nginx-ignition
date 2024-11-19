package br.com.dillmann.nginxignition.api.certificate.handler

import br.com.dillmann.nginxignition.api.certificate.model.CertificateConverter
import br.com.dillmann.nginxignition.api.certificate.model.IssueCertificateRequest
import br.com.dillmann.nginxignition.core.certificate.command.IssueCertificateCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.api.common.request.payload

internal class IssueCertificateHandler(
    private val issueCertificateCommand: IssueCertificateCommand,
    private val converter: CertificateConverter,
): RequestHandler {
    override suspend fun handle(call: ApiCall) {
        val requestPayload: IssueCertificateRequest
        try {
            requestPayload = call.payload()
        } catch (ex: Exception) {
            call.respond(HttpStatus.BAD_REQUEST, mapOf("message" to ex.message))
            return
        }

        val request = converter.toDomainModel(requestPayload)
        val issueOutput = issueCertificateCommand.issue(request)
        val responsePayload = converter.toResponse(issueOutput)
        val status = if (issueOutput.success) HttpStatus.OK else HttpStatus.BAD_REQUEST
        call.respond(status, responsePayload)
    }
}
