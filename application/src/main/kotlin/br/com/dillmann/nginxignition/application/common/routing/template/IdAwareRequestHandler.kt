package br.com.dillmann.nginxignition.application.common.routing.template

import br.com.dillmann.nginxignition.application.common.routing.RequestHandler
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import java.util.UUID

interface IdAwareRequestHandler: RequestHandler {
    override suspend fun handle(call: RoutingCall) {
        val id = runCatching { call.request.pathVariables["id"].let(UUID::fromString) }.getOrNull()
        if (id == null) {
            call.respond(HttpStatusCode.BadRequest)
            return
        }

        handle(call, id)
    }

    suspend fun handle(call: RoutingCall, id: UUID)
}
