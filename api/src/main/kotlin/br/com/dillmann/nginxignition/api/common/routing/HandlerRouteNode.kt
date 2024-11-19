package br.com.dillmann.nginxignition.api.common.routing

import br.com.dillmann.nginxignition.api.common.request.HttpVerb
import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler

data class HandlerRouteNode(val verb: HttpVerb, val path: String?, val handler: RequestHandler) : RouteNode
