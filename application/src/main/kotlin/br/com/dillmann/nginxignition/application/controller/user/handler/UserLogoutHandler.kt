package br.com.dillmann.nginxignition.application.controller.user.handler

import br.com.dillmann.nginxignition.application.common.rbac.RbacJwtFacade
import br.com.dillmann.nginxignition.application.common.routing.RequestHandler
import io.ktor.http.*
import io.ktor.server.auth.*
import io.ktor.server.auth.jwt.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

class UserLogoutHandler(private val jwtFacade: RbacJwtFacade): RequestHandler {
    override suspend fun handle(call: RoutingCall) {
        call.principal<JWTCredential>()?.let(jwtFacade::revokeCredentials)
        call.respond(HttpStatusCode.NoContent)
    }
}
