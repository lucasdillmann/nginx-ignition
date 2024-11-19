package br.com.dillmann.nginxignition.api.common.routing

interface RouteProvider {
    fun apiRoutes(): RouteNode
}
