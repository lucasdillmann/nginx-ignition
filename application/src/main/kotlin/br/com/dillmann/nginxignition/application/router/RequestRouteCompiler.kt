package br.com.dillmann.nginxignition.application.router

import br.com.dillmann.nginxignition.api.common.request.HttpMethod as Method
import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.api.common.routing.*
import br.com.dillmann.nginxignition.application.router.interceptor.AuthenticationRequiredInterceptor
import br.com.dillmann.nginxignition.application.router.interceptor.RoleRequiredInterceptor
import java.util.regex.Pattern

internal class RequestRouteCompiler(private val routeProviders: List<RouteProvider>) {
    data class CompiledRoute(
        val method: String,
        val pattern: Pattern,
        val handler: RequestHandler,
    )

    private data class Route(
        val method: Method,
        val path: String,
        val handler: RequestHandler,
    )

    fun compile(): List<CompiledRoute> =
        routeProviders
            .flatMap { it.apiRoutes().handlers() }
            .sortedByDescending { it.path.length }
            .map {
                CompiledRoute(
                    pattern = it.path.toRoutePattern(),
                    handler = it.handler,
                    method = it.method.name,
                )
            }

    private fun RouteNode.handlers(): List<Route> =
        when (this) {
            is PathPrefixRouteNode ->
                flatten().map { it.copy(path = "$path${it.path}") }
            is AuthenticationRequiredRouteNode ->
                flatten().map { it.copy(handler = AuthenticationRequiredInterceptor(it.handler)) }
            is RoleRequiredRouteNode ->
                flatten().map { it.copy(handler = RoleRequiredInterceptor(role, it.handler)) }
            is HandlerRouteNode ->
                listOf(Route(method, path ?: "", handler))
        }

    private fun CompositeRouteNode.flatten() =
        children().flatMap { it.handlers() }

    private fun String.toRoutePattern(): Pattern {
        val suffix = if (endsWith("/**")) "" else "\$"
        val pattern = replace("{", "(?<")
            .replace("}", ">[^/]+)")
            .removeSuffix("**")
        return Pattern.compile("^$pattern$suffix")
    }
}
