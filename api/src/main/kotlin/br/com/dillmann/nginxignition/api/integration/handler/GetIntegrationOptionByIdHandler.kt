package br.com.dillmann.nginxignition.api.integration.handler

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.IdAwareRequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond
import br.com.dillmann.nginxignition.api.integration.IntegrationConverter
import br.com.dillmann.nginxignition.core.integration.command.GetIntegrationOptionByIdCommand

internal class GetIntegrationOptionByIdHandler(
    private val command: GetIntegrationOptionByIdCommand,
    private val converter: IntegrationConverter,
): IdAwareRequestHandler {
    override suspend fun handle(call: ApiCall, id: String) {
        withIntegrationExceptionHandler(call) {
            val optionId = runCatching { call.pathVariables()["optionId"] }.getOrNull()
            if (optionId == null) {
                call.respond(HttpStatus.BAD_REQUEST)
                return@withIntegrationExceptionHandler
            }

            val payload = command.getIntegrationOptionById(id, optionId)?.let(converter::toResponse)
            if (payload == null)
                call.respond(HttpStatus.NOT_FOUND)
            else
                call.respond(HttpStatus.OK, payload)
        }
    }
}
