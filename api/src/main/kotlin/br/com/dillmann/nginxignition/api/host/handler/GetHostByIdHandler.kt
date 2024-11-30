package br.com.dillmann.nginxignition.api.host.handler

import br.com.dillmann.nginxignition.api.host.HostConverter
import br.com.dillmann.nginxignition.core.host.command.GetHostCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.UuidAwareRequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond
import java.util.UUID

internal class GetHostByIdHandler(
    private val getCommand: GetHostCommand,
    private val converter: HostConverter,
): UuidAwareRequestHandler {
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
