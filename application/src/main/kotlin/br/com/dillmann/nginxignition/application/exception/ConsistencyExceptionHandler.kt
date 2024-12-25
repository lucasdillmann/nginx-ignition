package br.com.dillmann.nginxignition.application.exception

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.respond
import br.com.dillmann.nginxignition.core.common.validation.ConsistencyException
import kotlinx.serialization.Serializable

// TODO: Reintegrate this
class ConsistencyExceptionHandler {
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

    suspend fun handle(call: ApiCall, ex: ConsistencyException) {
        val payload = Response(
            message = "One or more consistency problems were found",
            consistencyProblems = ex.violations.map { Response.Error(it.path, it.message) },
        )

        call.respond(HttpStatus.BAD_REQUEST, payload)
    }
}
