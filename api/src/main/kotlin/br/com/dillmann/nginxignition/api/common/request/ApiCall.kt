package br.com.dillmann.nginxignition.api.common.request

import br.com.dillmann.nginxignition.api.common.authorization.Subject
import java.io.InputStream
import kotlin.reflect.KClass
import kotlin.reflect.KType

interface ApiCall {
    suspend fun <T: Any> payload(contract: KClass<T>): T
    suspend fun <T: Any> respond(status: HttpStatus, payload: T, clazz: KClass<out T>, type: KType)
    suspend fun respond(status: HttpStatus)
    suspend fun respondRaw(status: HttpStatus, headers: Map<String, String>, payload: InputStream, payloadSize: Long)
    suspend fun headers(): Map<String, List<String>>
    suspend fun queryParams(): Map<String, String>
    suspend fun pathVariables(): Map<String, String>
    suspend fun principal(): Subject?
    suspend fun uri(): String
}
