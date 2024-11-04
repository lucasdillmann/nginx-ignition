package br.com.dillmann.nginxignition.application.common.configuration

import br.com.dillmann.nginxignition.application.common.rbac.RbacJwtFacade
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.auth.*
import io.ktor.server.auth.jwt.*
import io.ktor.server.response.*
import org.koin.ktor.ext.inject

fun Application.configureRbac() {
    val rbacJwt by inject<RbacJwtFacade>()

    install(Authentication) {
        jwt(RbacJwtFacade.UNIQUE_IDENTIFIER) {
            realm = RbacJwtFacade.UNIQUE_IDENTIFIER
            verifier(rbacJwt.buildVerifier())
            validate { rbacJwt.checkCredentials(it) }
            challenge { _, _ ->
                val payload = mapOf("message" to "Invalid or expired JWT token")
                call.respond(HttpStatusCode.Unauthorized, payload)
            }
        }
    }
}
