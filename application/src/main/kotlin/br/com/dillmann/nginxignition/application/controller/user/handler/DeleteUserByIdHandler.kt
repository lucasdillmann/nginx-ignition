package br.com.dillmann.nginxignition.application.controller.user.handler

import br.com.dillmann.nginxignition.application.common.routing.template.IdAwareRequestHandler
import br.com.dillmann.nginxignition.core.user.command.DeleteUserCommand
import io.ktor.http.*
import io.ktor.server.auth.*
import io.ktor.server.auth.jwt.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import java.util.*

class DeleteUserByIdHandler(
    private val deleteCommand: DeleteUserCommand,
): IdAwareRequestHandler {
    override suspend fun handle(call: RoutingCall, id: UUID) {
        val currentUser = call.principal<JWTPrincipal>()
        if (currentUser?.subject == id.toString()) {
            val payload = mapOf("message" to "You cannot delete your own user")
            call.respond(HttpStatusCode.BadRequest, payload)
            return
        }

        deleteCommand.deleteById(id)
        call.respond(HttpStatusCode.NoContent)
    }
}
