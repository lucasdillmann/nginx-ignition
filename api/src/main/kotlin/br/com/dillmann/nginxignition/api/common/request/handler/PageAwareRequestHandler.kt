package br.com.dillmann.nginxignition.api.common.request.handler

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus

internal interface PageAwareRequestHandler: RequestHandler {
    private companion object {
        private val ALLOWED_SIZE_RANGE = 1..1000
    }

    override suspend fun handle(call: ApiCall) {
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

    suspend fun handle(call: ApiCall, pageNumber: Int, pageSize: Int)

    private suspend fun ApiCall.sendErrorMessage(message: String) {
        respond(
            HttpStatus.BAD_REQUEST,
            mapOf("message" to message),
        )
    }

    private suspend fun ApiCall.queryParam(name: String) =
        runCatching { queryParams()[name]?.toInt() }.getOrNull()
}