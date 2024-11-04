package br.com.dillmann.nginxignition.application.controller.user.handler

import br.com.dillmann.nginxignition.application.controller.user.model.UserConverter
import br.com.dillmann.nginxignition.application.controller.user.model.UserRequest
import br.com.dillmann.nginxignition.core.user.command.SaveUserCommand
import io.ktor.http.*
import io.ktor.server.auth.*
import io.ktor.server.auth.jwt.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import java.util.*

class UpdateUserByIdHandler(
    private val saveCommand: SaveUserCommand,
    private val converter: UserConverter,
) {
    suspend fun handle(call: RoutingCall) {
        val userId = runCatching { call.request.pathVariables["id"].let(UUID::fromString) }.getOrNull()
        if (userId == null) {
            call.respond(HttpStatusCode.BadRequest)
            return
        }

        val payload = call.receive<UserRequest>()
        val currentUser = call.principal<JWTPrincipal>()
        if (currentUser?.subject == userId.toString() && !payload.enabled) {
            val responsePayload = mapOf("message" to "You cannot disable your own user")
            call.respond(HttpStatusCode.BadRequest, responsePayload)
            return
        }

        val user = converter.toDomainModel(payload).copy(id = userId)
        saveCommand.save(user)
        call.respond(HttpStatusCode.NoContent)
    }
}
