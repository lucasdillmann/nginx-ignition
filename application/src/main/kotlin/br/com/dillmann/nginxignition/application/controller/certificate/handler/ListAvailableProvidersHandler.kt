package br.com.dillmann.nginxignition.application.controller.certificate.handler

import br.com.dillmann.nginxignition.application.common.routing.RequestHandler
import br.com.dillmann.nginxignition.application.controller.certificate.model.CertificateConverter
import br.com.dillmann.nginxignition.core.certificate.command.GetAvailableProvidersCommand
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

class ListAvailableProvidersHandler(
    private val listProvidersCommand: GetAvailableProvidersCommand,
    private val converter: CertificateConverter,
): RequestHandler {
    override suspend fun handle(call: RoutingCall) {
        val providers = listProvidersCommand.getAvailableProviders()
        val payload = providers.map(converter::toResponse)
        call.respond(HttpStatusCode.OK, payload)
    }
}
