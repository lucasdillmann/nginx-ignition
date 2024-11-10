package br.com.dillmann.nginxignition.application.controller.certificate.handler

import br.com.dillmann.nginxignition.application.common.routing.RequestHandler
import br.com.dillmann.nginxignition.application.controller.certificate.model.CertificateConverter
import br.com.dillmann.nginxignition.application.controller.certificate.model.IssueCertificateRequest
import br.com.dillmann.nginxignition.core.certificate.command.IssueCertificateCommand
import io.ktor.http.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

class IssueCertificateHandler(
    private val issueCertificateCommand: IssueCertificateCommand,
    private val converter: CertificateConverter,
): RequestHandler {
    override suspend fun handle(call: RoutingCall) {
        val requestPayload: IssueCertificateRequest
        try {
            requestPayload = call.receive()
        } catch (ex: Exception) {
            call.respond(HttpStatusCode.BadRequest, mapOf("message" to ex.message))
            return
        }

        val request = converter.toDomainModel(requestPayload)
        val issueOutput = issueCertificateCommand.issue(request)
        val responsePayload = converter.toResponse(issueOutput)
        val status = if (issueOutput.success) HttpStatusCode.OK else HttpStatusCode.BadRequest
        call.respond(status, responsePayload)
    }
}
