package br.com.dillmann.nginxignition.api.common.logs

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.respond

private const val DEFAULT_LINE_COUNT = 50
private const val MAXIMUM_LINE_COUNT = 10_000

internal interface LogRequestHandler {
    suspend fun parseLinesAmount(call: ApiCall): Int? {
        val lineCount = runCatching { call.queryParams()["lines"]?.toInt() }.getOrNull() ?: DEFAULT_LINE_COUNT
        return if (lineCount in 1..MAXIMUM_LINE_COUNT) lineCount else null
    }

    suspend fun sendLinesAmountErrorResponse(call: ApiCall) {
        call.respond(
            HttpStatus.BAD_REQUEST,
            mapOf("message" to "Lines amount should be between 1 and $MAXIMUM_LINE_COUNT"),
        )
    }
}
