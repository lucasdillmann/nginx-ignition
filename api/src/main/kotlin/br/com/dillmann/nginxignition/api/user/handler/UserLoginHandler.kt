package br.com.dillmann.nginxignition.api.user.handler

import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.api.user.model.UserLoginRequest
import br.com.dillmann.nginxignition.api.user.model.UserLoginResponse
import br.com.dillmann.nginxignition.core.user.command.AuthenticateUserCommand
import br.com.dillmann.nginxignition.api.common.authorization.Authorizer
import br.com.dillmann.nginxignition.api.common.authorization.Subject
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.payload

internal class UserLoginHandler(
    private val authenticateCommand: AuthenticateUserCommand,
    private val authorizer: Authorizer,
): RequestHandler {
    override suspend fun handle(call: ApiCall) {
        val payload = call.payload<UserLoginRequest>()
        val authenticatedUser = authenticateCommand.authenticate(payload.username, payload.password)
        if (authenticatedUser == null) {
            call.respond(HttpStatus.FORBIDDEN)
            return
        }

        val subject = Subject(userId = authenticatedUser.id)
        val token = authorizer.buildToken(subject)
        val responsePayload = UserLoginResponse(token)
        call.respond(HttpStatus.OK, responsePayload)
    }
}
