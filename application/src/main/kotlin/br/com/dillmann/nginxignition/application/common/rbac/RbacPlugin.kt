package br.com.dillmann.nginxignition.application.common.rbac

import br.com.dillmann.nginxignition.core.user.User
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.auth.*
import io.ktor.server.auth.jwt.*
import io.ktor.server.response.*
import org.koin.java.KoinJavaComponent.getKoin

class RbacPluginConfiguration {
    var requiredRoles: Set<User.Role>? = null
}

val RbacPlugin = createRouteScopedPlugin(
    name = "RbacPlugin",
    createConfiguration = ::RbacPluginConfiguration,
) {
    on(AuthenticationChecked) { call ->
        val replyWithAccessDenied: suspend () -> Unit = {
            val payload = mapOf("message" to "Access denied")
            call.respond(HttpStatusCode.Unauthorized, payload)
        }

        val credentials = call.principal<JWTPrincipal>()
        if (credentials == null) {
            replyWithAccessDenied()
            return@on
        }

        val userRole = runCatching { credentials.payload.getClaim("role").`as`(User.Role::class.java) }.getOrNull()
        val requiredRoles = pluginConfig.requiredRoles
        val authorized = requiredRoles.isNullOrEmpty() || userRole != null && userRole in requiredRoles
        if (!authorized) {
            replyWithAccessDenied()
            return@on
        }

        val newToken = getKoin().get<RbacJwtFacade>().refreshToken(credentials)
        if (newToken != null) {
            call.response.header(HttpHeaders.Authorization, "Bearer $newToken")
        }
    }
}
