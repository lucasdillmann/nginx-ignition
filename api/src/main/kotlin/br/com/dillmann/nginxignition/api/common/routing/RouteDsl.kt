package br.com.dillmann.nginxignition.api.common.routing

import br.com.dillmann.nginxignition.core.user.User
import br.com.dillmann.nginxignition.api.common.request.HttpVerb
import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler

internal typealias RouteConfigurer = RouteNodeBuilder.() -> Unit

internal class RouteNodeBuilder(private val parent: CompositeRouteNode) {
    fun get(path: String, handler: RequestHandler) =
        configureHandlerNode(HttpVerb.GET, path, handler)

    fun get(handler: RequestHandler) =
        configureHandlerNode(HttpVerb.GET, null, handler)

    fun post(path: String, handler: RequestHandler) =
        configureHandlerNode(HttpVerb.POST, path, handler)

    fun post(handler: RequestHandler) =
        configureHandlerNode(HttpVerb.POST, null, handler)

    fun put(path: String, handler: RequestHandler) =
        configureHandlerNode(HttpVerb.PUT, path, handler)

    fun put(handler: RequestHandler) =
        configureHandlerNode(HttpVerb.PUT, null, handler)

    fun delete(path: String, handler: RequestHandler) =
        configureHandlerNode(HttpVerb.DELETE, path, handler)

    fun delete(handler: RequestHandler) =
        configureHandlerNode(HttpVerb.DELETE, null, handler)

    fun patch(path: String, handler: RequestHandler) =
        configureHandlerNode(HttpVerb.PATCH, path, handler)

    fun patch(handler: RequestHandler) =
        configureHandlerNode(HttpVerb.PATCH, null, handler)

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

    private fun configureHandlerNode(verb: HttpVerb, path: String?, handler: RequestHandler) {
        val route = HandlerRouteNode(verb, path, handler)
        parent.addChild(route)
    }
}

internal fun routes(path: String, customizer: RouteNodeBuilder.() -> Unit): RouteNode {
    val child = PathPrefixRouteNode(path)
    RouteNodeBuilder(child).customizer()
    return child
}
