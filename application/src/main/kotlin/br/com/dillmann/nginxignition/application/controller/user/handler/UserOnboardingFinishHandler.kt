package br.com.dillmann.nginxignition.application.controller.user.handler

import br.com.dillmann.nginxignition.application.common.rbac.RbacJwtFacade
import br.com.dillmann.nginxignition.application.common.routing.RequestHandler
import br.com.dillmann.nginxignition.application.controller.user.model.UserConverter
import br.com.dillmann.nginxignition.application.controller.user.model.UserLoginResponse
import br.com.dillmann.nginxignition.application.controller.user.model.UserRequest
import br.com.dillmann.nginxignition.core.user.User
import br.com.dillmann.nginxignition.core.user.command.AuthenticateUserCommand
import br.com.dillmann.nginxignition.core.user.command.GetUserCountCommand
import br.com.dillmann.nginxignition.core.user.command.SaveUserCommand
import io.ktor.http.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

class UserOnboardingFinishHandler(
    private val userCountCommand: GetUserCountCommand,
    private val saveUserCommand: SaveUserCommand,
    private val authenticateCommand: AuthenticateUserCommand,
    private val converter: UserConverter,
    private val jwtFacade: RbacJwtFacade,
): RequestHandler {
    override suspend fun handle(call: RoutingCall) {
        if (userCountCommand.count() > 0) {
            call.respond(HttpStatusCode.Forbidden)
            return
        }

        val userRequest = call
            .receive<UserRequest>()
            .let(converter::toDomainModel)
            .copy(
                enabled = true,
                role = User.Role.ADMIN,
            )

        saveUserCommand.save(userRequest)

        val user = authenticateCommand.authenticate(userRequest.username, userRequest.password!!)!!
        val token = jwtFacade.buildToken(user)
        val payload = UserLoginResponse(token)
        call.respond(HttpStatusCode.OK, payload)
    }
}
