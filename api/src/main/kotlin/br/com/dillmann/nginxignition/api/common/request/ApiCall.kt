package br.com.dillmann.nginxignition.api.common.request

import br.com.dillmann.nginxignition.api.common.authorization.Subject
import kotlin.reflect.KClass

interface ApiCall {
    suspend fun <T: Any> payload(contract: KClass<T>): T
    suspend fun respond(status: HttpStatus = HttpStatus.NO_CONTENT, payload: Any? = null)
    suspend fun headers(): Map<String, List<String>>
    suspend fun queryParams(): Map<String, String>
    suspend fun pathVariables(): Map<String, String>
    suspend fun principal(): Subject?
}
