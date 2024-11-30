package br.com.dillmann.nginxignition.api.integration.handler

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond
import br.com.dillmann.nginxignition.api.integration.IntegrationConverter
import br.com.dillmann.nginxignition.core.integration.command.ListIntegrationsCommand

internal class ListIntegrationsHandler(
    private val listCommand: ListIntegrationsCommand,
    private val converter: IntegrationConverter,
): RequestHandler {
    override suspend fun handle(call: ApiCall) {
        val payload = listCommand.getIntegrations().map(converter::toResponse)
        call.respond(HttpStatus.OK, payload)
    }
}
