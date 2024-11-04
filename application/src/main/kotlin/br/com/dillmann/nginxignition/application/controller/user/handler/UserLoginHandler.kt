package br.com.dillmann.nginxignition.application.controller.user.handler

import br.com.dillmann.nginxignition.application.common.rbac.RbacJwtFacade
import br.com.dillmann.nginxignition.application.controller.user.model.UserLoginRequest
import br.com.dillmann.nginxignition.application.controller.user.model.UserLoginResponse
import br.com.dillmann.nginxignition.core.user.command.AuthenticateUserCommand
import io.ktor.http.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

class UserLoginHandler(
    private val authenticateCommand: AuthenticateUserCommand,
    private val jwtFacade: RbacJwtFacade,
) {
    suspend fun handle(call: RoutingCall) {
        val payload = call.receive<UserLoginRequest>()
        val authenticatedUser = authenticateCommand.authenticate(payload.username, payload.password)
        if (authenticatedUser == null) {
            call.respond(HttpStatusCode.Forbidden)
            return
        }

        val token = jwtFacade.buildToken(authenticatedUser)
        val responsePayload = UserLoginResponse(token)
        call.respond(HttpStatusCode.OK, responsePayload)
    }
}
