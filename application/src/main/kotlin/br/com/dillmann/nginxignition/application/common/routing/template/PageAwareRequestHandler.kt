package br.com.dillmann.nginxignition.application.common.routing.template

import br.com.dillmann.nginxignition.application.common.routing.RequestHandler
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

interface PageAwareRequestHandler: RequestHandler {
    private companion object {
        private val ALLOWED_SIZE_RANGE = 1..1000
    }

    override suspend fun handle(call: RoutingCall) {
        val pageSize = call.queryParam("pageSize") ?: 25
        val pageNumber = call.queryParam("pageNumber") ?: 0

        if (pageSize !in ALLOWED_SIZE_RANGE) {
            call.sendErrorMessage(
                "Page size must be between ${ALLOWED_SIZE_RANGE.min()} and ${ALLOWED_SIZE_RANGE.max()}",
            )
            return
        }

        if (pageNumber < 0) {
            call.sendErrorMessage("Page number must be greater than or equal to 0")
            return
        }

        handle(call, pageNumber, pageSize)
    }

    suspend fun handle(call: RoutingCall, pageNumber: Int, pageSize: Int)

    private suspend fun RoutingCall.sendErrorMessage(message: String) {
        respond(
            HttpStatusCode.BadRequest,
            mapOf("message" to message),
        )
    }

    private fun RoutingCall.queryParam(name: String) =
        runCatching { queryParameters[name]?.toInt() }.getOrNull()
}
