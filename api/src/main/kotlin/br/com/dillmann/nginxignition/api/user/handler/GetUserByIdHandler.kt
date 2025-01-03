package br.com.dillmann.nginxignition.api.user.handler

import br.com.dillmann.nginxignition.api.user.UserConverter
import br.com.dillmann.nginxignition.core.user.command.GetUserCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.UuidAwareRequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond
import java.util.*

internal class GetUserByIdHandler(
    private val getCommand: GetUserCommand,
    private val converter: UserConverter,
): UuidAwareRequestHandler {
    override suspend fun handle(call: ApiCall, id: UUID) {
        val user = getCommand.getById(id)
        if (user != null) {
            val payload = converter.toResponse(user)
            call.respond(HttpStatus.OK, payload)
        } else {
            call.respond(HttpStatus.NOT_FOUND)
        }
    }
}
