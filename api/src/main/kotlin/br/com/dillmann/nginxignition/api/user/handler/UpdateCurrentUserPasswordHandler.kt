package br.com.dillmann.nginxignition.api.user.handler

import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.payload
import br.com.dillmann.nginxignition.api.user.model.UserPasswordUpdateRequest
import br.com.dillmann.nginxignition.core.user.command.UpdateUserPasswordCommand

internal class UpdateCurrentUserPasswordHandler(
    private val updatePasswordCommand: UpdateUserPasswordCommand,
): RequestHandler {
    override suspend fun handle(call: ApiCall) {
        val payload = call.payload<UserPasswordUpdateRequest>()
        val principal = call.principal()!!
        updatePasswordCommand.updatePassword(principal.userId, payload.currentPassword, payload.newPassword)
        call.respond(HttpStatus.NO_CONTENT)
    }
}
