package br.com.dillmann.nginxignition.api.common.request.handler

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.respond

internal interface PageAwareRequestHandler: RequestHandler {
    private companion object {
        private val ALLOWED_SIZE_RANGE = 1..1000
        private const val DEFAULT_PAGE_SIZE = 25
        private const val DEFAULT_PAGE_NUMBER = 0
    }

    override suspend fun handle(call: ApiCall) {
        val pageSize = call.queryParam("pageSize") ?: DEFAULT_PAGE_SIZE
        val pageNumber = call.queryParam("pageNumber") ?: DEFAULT_PAGE_NUMBER
        val searchTerms = call.queryParams()["searchTerms"]

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

        handle(call, pageNumber, pageSize, searchTerms)
    }

    suspend fun handle(call: ApiCall, pageNumber: Int, pageSize: Int, searchTerms: String?)

    private suspend fun ApiCall.sendErrorMessage(message: String) {
        respond(
            HttpStatus.BAD_REQUEST,
            mapOf("message" to message),
        )
    }

    private suspend fun ApiCall.queryParam(name: String) =
        runCatching { queryParams()[name]?.toInt() }.getOrNull()
}
