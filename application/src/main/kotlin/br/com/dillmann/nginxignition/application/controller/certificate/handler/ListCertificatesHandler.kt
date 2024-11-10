package br.com.dillmann.nginxignition.application.controller.certificate.handler

import br.com.dillmann.nginxignition.application.common.routing.RequestHandler
import br.com.dillmann.nginxignition.application.controller.certificate.model.CertificateConverter
import br.com.dillmann.nginxignition.core.certificate.command.ListCertificateCommand
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

class ListCertificatesHandler(
    private val listCommand: ListCertificateCommand,
    private val converter: CertificateConverter,
): RequestHandler {
    override suspend fun handle(call: RoutingCall) {
        val pageSize = runCatching { call.request.queryParameters["pageSize"]?.toInt() }.getOrNull() ?: 10
        val pageNumber = runCatching { call.request.queryParameters["pageNumber"]?.toInt() }.getOrNull() ?: 0

        val page = listCommand.list(pageSize, pageNumber)
        val payload = converter.toResponse(page)
        call.respond(HttpStatusCode.OK, payload)
    }
}
