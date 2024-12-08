package br.com.dillmann.nginxignition.application.configuration

import br.com.dillmann.nginxignition.application.rbac.RbacJwtFacade
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
            challenge { _, _ -> call.respond(HttpStatusCode.Unauthorized) }
        }
    }
}
