package br.com.dillmann.nginxignition.application.rbac

import br.com.dillmann.nginxignition.core.user.User
import io.ktor.server.auth.*
import io.ktor.server.routing.*
import io.ktor.utils.io.*

@KtorDsl
fun Route.requireRole(vararg roles: User.Role, configuration: Route.() -> Unit) =
    withRbac(configuration) {
        requiredRoles = roles.toSet()
    }

@KtorDsl
fun Route.requireAuthentication(configuration: Route.() -> Unit) =
    withRbac(configuration)

private fun Route.withRbac(
    routeConfiguration: Route.() -> Unit,
    rbacConfiguration: RbacPluginConfiguration.() -> Unit = {},
) {
    authenticate(RbacJwtFacade.UNIQUE_IDENTIFIER) {
        install(RbacPlugin, rbacConfiguration)
        routeConfiguration()
    }
}
