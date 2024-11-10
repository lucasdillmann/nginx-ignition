package br.com.dillmann.nginxignition.application.controller.host.handler

import br.com.dillmann.nginxignition.application.common.routing.template.IdAwareRequestHandler
import br.com.dillmann.nginxignition.core.host.command.DeleteHostCommand
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import java.util.UUID

class DeleteHostByIdHandler(
    private val deleteCommand: DeleteHostCommand,
): IdAwareRequestHandler {
    override suspend fun handle(call: RoutingCall, id: UUID) {
        deleteCommand.deleteById(id)
        call.respond(HttpStatusCode.NoContent)
    }
}
