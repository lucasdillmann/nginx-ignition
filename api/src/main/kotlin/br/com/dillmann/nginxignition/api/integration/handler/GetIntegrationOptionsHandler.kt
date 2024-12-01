package br.com.dillmann.nginxignition.api.integration.handler

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.PageAwareRequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond
import br.com.dillmann.nginxignition.api.integration.IntegrationConverter
import br.com.dillmann.nginxignition.core.integration.command.ListIntegrationOptionsCommand

internal class GetIntegrationOptionsHandler(
    private val getOptionsCommand: ListIntegrationOptionsCommand,
    private val converter: IntegrationConverter,
): PageAwareRequestHandler {
    override suspend fun handle(call: ApiCall, pageNumber: Int, pageSize: Int) {
        withIntegrationExceptionHandler(call) {
            val id = runCatching { call.pathVariables()["id"] }.getOrNull()
            if (id == null) {
                call.respond(HttpStatus.BAD_REQUEST)
                return@withIntegrationExceptionHandler
            }

            val page = getOptionsCommand.getIntegrationOptions(id, pageNumber, pageSize)
            val payload = converter.toResponse(page)
            call.respond(HttpStatus.OK, payload)
        }
    }
}
