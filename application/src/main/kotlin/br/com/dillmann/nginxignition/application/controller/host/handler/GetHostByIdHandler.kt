package br.com.dillmann.nginxignition.application.controller.host.handler

import br.com.dillmann.nginxignition.application.common.routing.template.IdAwareRequestHandler
import br.com.dillmann.nginxignition.application.controller.host.model.HostConverter
import br.com.dillmann.nginxignition.core.host.command.GetHostCommand
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import java.util.UUID

class GetHostByIdHandler(
    private val getCommand: GetHostCommand,
    private val converter: HostConverter,
): IdAwareRequestHandler {
    override suspend fun handle(call: RoutingCall, id: UUID) {
        val host = getCommand.getById(id)
        if (host != null) {
            val payload = converter.toResponse(host)
            call.respond(HttpStatusCode.OK, payload)
        } else {
            call.respond(HttpStatusCode.NotFound)
        }
    }
}
