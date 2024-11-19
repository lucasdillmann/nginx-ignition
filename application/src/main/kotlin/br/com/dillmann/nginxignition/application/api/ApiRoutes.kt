package br.com.dillmann.nginxignition.application.api

import br.com.dillmann.nginxignition.application.rbac.requireAuthentication
import br.com.dillmann.nginxignition.application.rbac.requireRole
import br.com.dillmann.nginxignition.api.common.routing.*
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.routing.*
import org.koin.ktor.ext.getKoin

fun Application.apiRoutes() {
    val routeProviders = getKoin().getAll<RouteProvider>()

    routing {
        routeProviders.forEach {
            install(it.apiRoutes())
        }
    }
}

private fun Route.install(node: RouteNode) {
    when (node) {
        is AuthenticationRequiredRouteNode ->
            requireAuthentication {
                node.children().forEach { install(it) }
            }

        is RoleRequiredRouteNode ->
            requireRole(node.role) {
                node.children().forEach { install(it) }
            }

        is PathPrefixRouteNode ->
            route(node.path) {
                node.children().forEach { install(it) }
            }

        is HandlerRouteNode ->
            route(
                path = node.path ?: "",
                method = HttpMethod.parse(node.verb.name),
            ) {
                handle {
                    val callAdapter = KtorApiCallAdapter(call)
                    node.handler.handle(callAdapter)
                }
            }
    }
}
