package br.com.dillmann.nginxsidewheel.application.controller.host.handler

import br.com.dillmann.nginxsidewheel.application.controller.host.model.HostConverter
import br.com.dillmann.nginxsidewheel.application.controller.host.model.HostRequest
import br.com.dillmann.nginxsidewheel.core.host.command.SaveHostCommand
import io.ktor.http.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import java.util.UUID

class UpdateHostByIdHandler(
    private val saveCommand: SaveHostCommand,
    private val converter: HostConverter,
) {
    suspend fun handle(call: RoutingCall) {
        val hostId = runCatching { call.request.pathVariables["id"].let(UUID::fromString) }.getOrNull()
        if (hostId == null) {
            call.respond(HttpStatusCode.BadRequest)
            return
        }

        val payload = call.receive<HostRequest>()
        val host = converter.toDomainModel(payload).copy(id = hostId)
        saveCommand.save(host)
        call.respond(HttpStatusCode.NoContent)
    }
}
