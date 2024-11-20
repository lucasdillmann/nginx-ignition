package br.com.dillmann.nginxignition.api.common.routing

import br.com.dillmann.nginxignition.api.common.request.HttpMethod
import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler

data class HandlerRouteNode(val method: HttpMethod, val path: String?, val handler: RequestHandler) : RouteNode
