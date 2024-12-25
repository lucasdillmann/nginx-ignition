package br.com.dillmann.nginxignition.application.router.exception

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.respond
import br.com.dillmann.nginxignition.core.common.validation.ConsistencyException
import kotlinx.serialization.Serializable

internal class ExceptionHandler {
    @Serializable
    private data class Response(
        val message: String,
        val consistencyProblems: List<Error>,
    ) {
        @Serializable
        data class Error(
            val fieldPath: String,
            val message: String,
        )
    }

    suspend fun handle(call: ApiCall, ex: Throwable) {
        when (ex) {
            is ConsistencyException -> handle(call, ex)
            else -> call.respond(HttpStatus.INTERNAL_SERVER_ERROR)
        }
    }

    private suspend fun handle(call: ApiCall, ex: ConsistencyException) {
        val payload = Response(
            message = "One or more consistency problems were found",
            consistencyProblems = ex.violations.map { Response.Error(it.path, it.message) },
        )

        call.respond(HttpStatus.BAD_REQUEST, payload)
    }
}
