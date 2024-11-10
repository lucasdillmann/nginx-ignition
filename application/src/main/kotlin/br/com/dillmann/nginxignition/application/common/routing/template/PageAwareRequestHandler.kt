package br.com.dillmann.nginxignition.application.common.routing.template

import br.com.dillmann.nginxignition.application.common.routing.RequestHandler
import io.ktor.server.routing.*

interface PageAwareRequestHandler: RequestHandler {
    override suspend fun handle(call: RoutingCall) {
        val pageSize = runCatching { call.request.queryParameters["pageSize"]?.toInt() }.getOrNull() ?: 10
        val pageNumber = runCatching { call.request.queryParameters["pageNumber"]?.toInt() }.getOrNull() ?: 0

        handle(call, pageNumber, pageSize)
    }

    suspend fun handle(call: RoutingCall, pageNumber: Int, pageSize: Int)
}
