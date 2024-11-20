package br.com.dillmann.nginxignition.api.certificate.handler

import br.com.dillmann.nginxignition.api.certificate.model.CertificateConverter
import br.com.dillmann.nginxignition.core.certificate.command.ListCertificateCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.PageAwareRequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond

internal class ListCertificatesHandler(
    private val listCommand: ListCertificateCommand,
    private val converter: CertificateConverter,
): PageAwareRequestHandler {
    override suspend fun handle(call: ApiCall, pageNumber: Int, pageSize: Int) {
        val page = listCommand.list(pageSize, pageNumber)
        val payload = converter.toResponse(page)
        call.respond(HttpStatus.OK, payload)
    }
}
