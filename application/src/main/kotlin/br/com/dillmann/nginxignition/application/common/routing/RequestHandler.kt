package br.com.dillmann.nginxignition.application.common.routing

import io.ktor.server.routing.*

interface RequestHandler {
    suspend fun handle(call: RoutingCall)
}
