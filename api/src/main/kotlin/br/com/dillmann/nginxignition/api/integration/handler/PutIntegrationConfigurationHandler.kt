package br.com.dillmann.nginxignition.api.integration.handler

import br.com.dillmann.nginxignition.api.common.jsonobject.toUnwrappedMap
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.IdAwareRequestHandler
import br.com.dillmann.nginxignition.api.common.request.payload
import br.com.dillmann.nginxignition.api.common.request.respond
import br.com.dillmann.nginxignition.api.integration.model.IntegrationConfigurationRequest
import br.com.dillmann.nginxignition.core.integration.command.ConfigureIntegrationByIdCommand

internal class PutIntegrationConfigurationHandler(
    private val configureCommand: ConfigureIntegrationByIdCommand,
): IdAwareRequestHandler {
    override suspend fun handle(call: ApiCall, id: String) {
        withIntegrationExceptionHandler(call) {
            val payload = call.payload<IntegrationConfigurationRequest>()
            configureCommand.configureIntegration(id, payload.enabled, payload.parameters.toUnwrappedMap())
            call.respond(HttpStatus.NO_CONTENT, payload)
        }
    }
}
