package br.com.dillmann.nginxignition.api.common.routing

fun interface RouteProvider {
    fun apiRoutes(): RouteNode
}
