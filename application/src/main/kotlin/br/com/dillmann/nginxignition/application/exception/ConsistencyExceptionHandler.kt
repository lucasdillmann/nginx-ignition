package br.com.dillmann.nginxignition.application.exception

import br.com.dillmann.nginxignition.core.common.validation.ConsistencyException
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.response.*
import kotlinx.serialization.Serializable

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

    suspend fun handle(call: ApplicationCall, ex: ConsistencyException) {
        val payload = Response(
            message = "One or more consistency problems were found",
            consistencyProblems = ex.violations.map { Response.Error(it.path, it.message) },
        )

        call.respond(HttpStatusCode.BadRequest, payload)
    }
}
