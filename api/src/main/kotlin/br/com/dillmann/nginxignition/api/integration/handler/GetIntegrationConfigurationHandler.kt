package br.com.dillmann.nginxignition.api.integration.handler

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.IdAwareRequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond
import br.com.dillmann.nginxignition.api.integration.IntegrationConverter
import br.com.dillmann.nginxignition.core.integration.command.GetIntegrationByIdCommand

internal class GetIntegrationConfigurationHandler(
    private val getCommand: GetIntegrationByIdCommand,
    private val converter: IntegrationConverter,
): IdAwareRequestHandler {
    override suspend fun handle(call: ApiCall, id: String) {
        val payload = getCommand.getIntegrationById(id)?.let(converter::toResponse)

        if (payload == null)
            call.respond(HttpStatus.NOT_FOUND)
        else
            call.respond(HttpStatus.OK, payload)
    }
}
