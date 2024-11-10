package br.com.dillmann.nginxignition.application.controller.user.handler

import br.com.dillmann.nginxignition.application.common.routing.template.IdAwareRequestHandler
import br.com.dillmann.nginxignition.application.controller.user.model.UserConverter
import br.com.dillmann.nginxignition.core.user.command.GetUserCommand
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import java.util.*

class GetUserByIdHandler(
    private val getCommand: GetUserCommand,
    private val converter: UserConverter,
): IdAwareRequestHandler {
    override suspend fun handle(call: RoutingCall, id: UUID) {
        val user = getCommand.getById(id)
        if (user != null) {
            val payload = converter.toResponse(user)
            call.respond(HttpStatusCode.OK, payload)
        } else {
            call.respond(HttpStatusCode.NotFound)
        }
    }
}
