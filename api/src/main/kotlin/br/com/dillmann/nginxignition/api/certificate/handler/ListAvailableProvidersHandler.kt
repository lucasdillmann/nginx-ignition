package br.com.dillmann.nginxignition.api.certificate.handler

import br.com.dillmann.nginxignition.api.certificate.CertificateConverter
import br.com.dillmann.nginxignition.core.certificate.command.GetAvailableProvidersCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond

internal class ListAvailableProvidersHandler(
    private val listProvidersCommand: GetAvailableProvidersCommand,
    private val converter: CertificateConverter,
): RequestHandler {
    override suspend fun handle(call: ApiCall) {
        val providers = listProvidersCommand.getAvailableProviders()
        val payload = providers.map(converter::toResponse)
        call.respond(HttpStatus.OK, payload)
    }
}
