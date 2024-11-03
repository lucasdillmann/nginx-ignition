package br.com.dillmann.nginxignition.application.controller.host.handler

import br.com.dillmann.nginxignition.application.controller.host.model.HostConverter
import br.com.dillmann.nginxignition.core.host.command.GetHostCommand
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import java.util.UUID

class GetHostByIdHandler(
    private val getCommand: GetHostCommand,
    private val converter: HostConverter,
) {
    suspend fun handle(call: RoutingCall) {
        val hostId = runCatching { call.request.pathVariables["id"].let(UUID::fromString) }.getOrNull()
        if (hostId == null) {
            call.respond(HttpStatusCode.BadRequest)
            return
        }

        val host = getCommand.getById(hostId)
        if (host != null) {
            val payload = converter.toResponse(host)
            call.respond(HttpStatusCode.OK, payload)
        } else {
            call.respond(HttpStatusCode.NotFound)
        }
    }
}
