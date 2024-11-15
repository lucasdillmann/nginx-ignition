package br.com.dillmann.nginxignition.application.controller.user.handler

import br.com.dillmann.nginxignition.application.common.routing.RequestHandler
import br.com.dillmann.nginxignition.application.controller.user.model.UserConverter
import br.com.dillmann.nginxignition.core.user.command.GetUserCommand
import io.ktor.http.*
import io.ktor.server.auth.*
import io.ktor.server.auth.jwt.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import java.util.*

class CurrentUserHandler(
    private val getUserCommand: GetUserCommand,
    private val converter: UserConverter,
): RequestHandler {
    override suspend fun handle(call: RoutingCall) {
        val principal = call.principal<JWTPrincipal>()!!
        val userId = principal.payload.subject.let(UUID::fromString)
        val user = getUserCommand.getById(userId)!!
        val payload = converter.toResponse(user)
        call.respond(HttpStatusCode.OK, payload)
    }
}
