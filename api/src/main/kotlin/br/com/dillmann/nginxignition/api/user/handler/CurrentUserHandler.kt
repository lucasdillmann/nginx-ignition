package br.com.dillmann.nginxignition.api.user.handler

import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.api.user.model.UserConverter
import br.com.dillmann.nginxignition.core.user.command.GetUserCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus

internal class CurrentUserHandler(
    private val getUserCommand: GetUserCommand,
    private val converter: UserConverter,
): RequestHandler {
    override suspend fun handle(call: ApiCall) {
        val principal = call.principal()!!
        val user = getUserCommand.getById(principal.userId)!!
        val payload = converter.toResponse(user)
        call.respond(HttpStatus.OK, payload)
    }
}
