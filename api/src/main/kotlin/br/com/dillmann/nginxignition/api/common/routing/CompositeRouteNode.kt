package br.com.dillmann.nginxignition.api.common.routing

sealed class CompositeRouteNode: RouteNode {
    private val children = mutableListOf<RouteNode>()

    fun addChild(node: RouteNode) {
        children += node
    }

    fun children(): List<RouteNode> = children
}
