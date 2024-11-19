package br.com.dillmann.nginxignition.api.user.handler

import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.api.user.model.UserConverter
import br.com.dillmann.nginxignition.api.user.model.UserLoginResponse
import br.com.dillmann.nginxignition.api.user.model.UserRequest
import br.com.dillmann.nginxignition.core.user.User
import br.com.dillmann.nginxignition.core.user.command.AuthenticateUserCommand
import br.com.dillmann.nginxignition.core.user.command.GetUserCountCommand
import br.com.dillmann.nginxignition.core.user.command.SaveUserCommand
import br.com.dillmann.nginxignition.api.common.authorization.Authorizer
import br.com.dillmann.nginxignition.api.common.authorization.Subject
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.payload

internal class UserOnboardingFinishHandler(
    private val userCountCommand: GetUserCountCommand,
    private val saveUserCommand: SaveUserCommand,
    private val authenticateCommand: AuthenticateUserCommand,
    private val converter: UserConverter,
    private val authorizer: Authorizer,
): RequestHandler {
    override suspend fun handle(call: ApiCall) {
        if (userCountCommand.count() > 0) {
            call.respond(HttpStatus.FORBIDDEN)
            return
        }

        val userRequest = call
            .payload<UserRequest>()
            .let(converter::toDomainModel)
            .copy(
                enabled = true,
                role = User.Role.ADMIN,
            )

        saveUserCommand.save(userRequest)

        val user = authenticateCommand.authenticate(userRequest.username, userRequest.password!!)!!
        val subject = Subject(userId = user.id)
        val token = authorizer.buildToken(subject)
        val payload = UserLoginResponse(token)
        call.respond(HttpStatus.OK, payload)
    }
}
