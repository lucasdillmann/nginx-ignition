package br.com.dillmann.nginxignition.api.user.handler

import br.com.dillmann.nginxignition.api.user.model.UserConverter
import br.com.dillmann.nginxignition.api.user.model.UserRequest
import br.com.dillmann.nginxignition.core.user.command.SaveUserCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.IdAwareRequestHandler
import br.com.dillmann.nginxignition.api.common.request.payload
import java.util.*

internal class UpdateUserByIdHandler(
    private val saveCommand: SaveUserCommand,
    private val converter: UserConverter,
): IdAwareRequestHandler {
    override suspend fun handle(call: ApiCall, id: UUID) {
        val payload = call.payload<UserRequest>()
        val currentUser = call.principal()
        if (currentUser?.userId == id && !payload.enabled) {
            val responsePayload = mapOf("message" to "You cannot disable your own user")
            call.respond(HttpStatus.BAD_REQUEST, responsePayload)
            return
        }

        val user = converter.toDomainModel(payload).copy(id = id)
        saveCommand.save(user)
        call.respond(HttpStatus.NO_CONTENT)
    }
}
