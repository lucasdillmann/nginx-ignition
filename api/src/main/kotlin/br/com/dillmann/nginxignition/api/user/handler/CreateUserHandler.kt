package br.com.dillmann.nginxignition.api.user.handler

import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.api.user.model.UserConverter
import br.com.dillmann.nginxignition.api.user.model.UserRequest
import br.com.dillmann.nginxignition.core.user.command.SaveUserCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.payload
import java.util.*

internal class CreateUserHandler(
    private val saveCommand: SaveUserCommand,
    private val converter: UserConverter,
): RequestHandler {
    override suspend fun handle(call: ApiCall) {
        val payload = call.payload<UserRequest>()
        val user = converter.toDomainModel(payload).copy(id = UUID.randomUUID())
        val currentUser = call.principal()
        saveCommand.save(user, currentUser?.userId)
        call.respond(HttpStatus.NO_CONTENT)
    }
}
