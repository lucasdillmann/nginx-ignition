package br.com.dillmann.nginxignition.application.controller.host.handler

import br.com.dillmann.nginxignition.core.host.command.DeleteHostCommand
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import java.util.UUID

class DeleteHostByIdHandler(
    private val deleteCommand: DeleteHostCommand,
) {
    suspend fun handle(call: RoutingCall) {
        val hostId = runCatching { call.request.pathVariables["id"].let(UUID::fromString) }.getOrNull()
        if (hostId == null) {
            call.respond(HttpStatusCode.BadRequest)
            return
        }

        deleteCommand.deleteById(hostId)
        call.respond(HttpStatusCode.NoContent)
    }
}
