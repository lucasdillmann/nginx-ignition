package br.com.dillmann.nginxignition.application.controller.user.handler

import br.com.dillmann.nginxignition.application.common.routing.RequestHandler
import br.com.dillmann.nginxignition.application.controller.user.model.UserConverter
import br.com.dillmann.nginxignition.application.controller.user.model.UserRequest
import br.com.dillmann.nginxignition.core.user.command.SaveUserCommand
import io.ktor.http.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import java.util.*

class CreateUserHandler(
    private val saveCommand: SaveUserCommand,
    private val converter: UserConverter,
): RequestHandler {
    override suspend fun handle(call: RoutingCall) {
        val payload = call.receive<UserRequest>()
        val user = converter.toDomainModel(payload).copy(id = UUID.randomUUID())
        saveCommand.save(user)
        call.respond(HttpStatusCode.NoContent)
    }
}
