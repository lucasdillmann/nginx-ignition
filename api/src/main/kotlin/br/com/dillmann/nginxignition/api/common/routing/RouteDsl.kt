package br.com.dillmann.nginxignition.api.common.routing

import br.com.dillmann.nginxignition.core.user.User
import br.com.dillmann.nginxignition.api.common.request.HttpMethod
import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler

typealias RouteConfigurer = RouteNodeBuilder.() -> Unit

class RouteNodeBuilder(private val parent: CompositeRouteNode) {
    fun get(path: String, handler: RequestHandler) =
        configureHandlerNode(HttpMethod.GET, path, handler)

    fun get(handler: RequestHandler) =
        configureHandlerNode(HttpMethod.GET, null, handler)

    fun post(path: String, handler: RequestHandler) =
        configureHandlerNode(HttpMethod.POST, path, handler)

    fun post(handler: RequestHandler) =
        configureHandlerNode(HttpMethod.POST, null, handler)

    fun put(path: String, handler: RequestHandler) =
        configureHandlerNode(HttpMethod.PUT, path, handler)

    fun put(handler: RequestHandler) =
        configureHandlerNode(HttpMethod.PUT, null, handler)

    fun delete(path: String, handler: RequestHandler) =
        configureHandlerNode(HttpMethod.DELETE, path, handler)

    fun delete(handler: RequestHandler) =
        configureHandlerNode(HttpMethod.DELETE, null, handler)

    fun patch(path: String, handler: RequestHandler) =
        configureHandlerNode(HttpMethod.PATCH, path, handler)

    fun patch(handler: RequestHandler) =
        configureHandlerNode(HttpMethod.PATCH, null, handler)

    fun requireAuthentication(configurer: RouteConfigurer) =
        configureCompositeNode(AuthenticationRequiredRouteNode(), configurer)

    fun requireRole(role: User.Role, configurer: RouteConfigurer) =
        configureCompositeNode(RoleRequiredRouteNode(role), configurer)

    fun path(path: String, configurer: RouteConfigurer) =
        configureCompositeNode(PathPrefixRouteNode(path), configurer)

    private fun configureCompositeNode(child: CompositeRouteNode, configurer: RouteConfigurer) {
        RouteNodeBuilder(child).configurer()
        parent.addChild(child)
    }

    private fun configureHandlerNode(method: HttpMethod, path: String?, handler: RequestHandler) {
        val route = HandlerRouteNode(method, path, handler)
        parent.addChild(route)
    }
}

fun basePath(path: String, customizer: RouteNodeBuilder.() -> Unit): RouteNode {
    val child = PathPrefixRouteNode(path)
    RouteNodeBuilder(child).customizer()
    return child
}
