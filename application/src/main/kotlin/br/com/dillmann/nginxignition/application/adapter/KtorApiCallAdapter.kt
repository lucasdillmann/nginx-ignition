package br.com.dillmann.nginxignition.application.adapter

import br.com.dillmann.nginxignition.api.common.authorization.Subject
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import io.ktor.http.*
import io.ktor.server.auth.*
import io.ktor.server.auth.jwt.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.util.*
import java.util.UUID
import kotlin.reflect.KClass

class KtorApiCallAdapter(private val call: RoutingCall): ApiCall {
    override suspend fun <T : Any> payload(contract: KClass<T>): T =
        call.receive(contract)

    override suspend fun respond(status: HttpStatus, payload: Any?) {
        val httpStatus = HttpStatusCode.fromValue(status.code)
        if (payload == null) call.respond(httpStatus)
        else call.respond(httpStatus, payload)
    }

    override suspend fun headers(): Map<String, List<String>> =
        call.request.headers.toMap()

    override suspend fun queryParams(): Map<String, String> =
        call.request.queryParameters.toMap().mapValues { it.value.first() }

    override suspend fun pathVariables(): Map<String, String> =
        call.request.pathVariables.toMap().mapValues { it.value.first() }

    override suspend fun principal(): Subject? {
        val principal = call.principal<JWTPrincipal>()?.payload ?: return null
        val userId = principal.subject.let(UUID::fromString)
        return Subject(userId = userId)
    }
}
