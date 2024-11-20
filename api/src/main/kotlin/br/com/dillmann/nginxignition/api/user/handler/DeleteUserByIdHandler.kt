package br.com.dillmann.nginxignition.api.user.handler

import br.com.dillmann.nginxignition.api.common.request.handler.IdAwareRequestHandler
import br.com.dillmann.nginxignition.core.user.command.DeleteUserCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.respond
import java.util.*

internal class DeleteUserByIdHandler(
    private val deleteCommand: DeleteUserCommand,
): IdAwareRequestHandler {
    override suspend fun handle(call: ApiCall, id: UUID) {
        val currentUser = call.principal()
        if (currentUser?.userId == id) {
            val payload = mapOf("message" to "You cannot delete your own user")
            call.respond(HttpStatus.BAD_REQUEST, payload)
            return
        }

        deleteCommand.deleteById(id)
        call.respond(HttpStatus.NO_CONTENT)
    }
}
