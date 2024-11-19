package br.com.dillmann.nginxignition.api.host.handler

import br.com.dillmann.nginxignition.api.host.model.HostConverter
import br.com.dillmann.nginxignition.core.host.command.GetHostCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.IdAwareRequestHandler
import java.util.UUID

internal class GetHostByIdHandler(
    private val getCommand: GetHostCommand,
    private val converter: HostConverter,
): IdAwareRequestHandler {
    override suspend fun handle(call: ApiCall, id: UUID) {
        val host = getCommand.getById(id)
        if (host != null) {
            val payload = converter.toResponse(host)
            call.respond(HttpStatus.OK, payload)
        } else {
            call.respond(HttpStatus.NOT_FOUND)
        }
    }
}
