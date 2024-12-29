package br.com.dillmann.nginxignition.application.router.adapter

import br.com.dillmann.nginxignition.api.common.authorization.Subject
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.application.router.ResponseInterceptor
import com.sun.net.httpserver.HttpExchange
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import kotlinx.serialization.ExperimentalSerializationApi
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.decodeFromStream
import kotlinx.serialization.serializer
import org.apache.commons.io.IOUtils
import java.io.InputStream
import java.net.URLDecoder
import java.util.StringTokenizer
import kotlin.reflect.KClass
import kotlin.reflect.KType

@OptIn(ExperimentalSerializationApi::class)
internal data class HttpApiCallAdapter(
    private val exchange: HttpExchange,
    private val interceptors: List<ResponseInterceptor>,
    private val principal: Subject? = null,
    private val pathVariables: Map<String, String> = emptyMap(),
): ApiCall {
    @Suppress("UNCHECKED_CAST")
    override suspend fun <T : Any> payload(contract: KClass<T>): T =
        Json.decodeFromStream(
            Json.serializersModule.serializer(contract.javaObjectType),
            exchange.requestBody,
        ) as T

    override suspend fun <T : Any> respond(status: HttpStatus, payload: T, clazz: KClass<out T>, type: KType) {
        invokeInterceptors()

        val json = Json.encodeToString(serializer(type), payload).byteInputStream()
        exchange.responseHeaders["content-type"] = "application/json"
        exchange.sendResponseHeaders(status.code, json.available().toLong())

        withContext(Dispatchers.IO) {
            IOUtils.copy(json, exchange.responseBody)
        }
    }

    override suspend fun respond(status: HttpStatus) {
        invokeInterceptors()
        exchange.sendResponseHeaders(status.code, -1)
    }

    override suspend fun respondRaw(
        status: HttpStatus,
        headers: Map<String, String>,
        payload: InputStream,
        payloadSize: Long,
    ) {
        invokeInterceptors()
        headers.forEach { (key, value) -> exchange.responseHeaders[key] = value }
        exchange.sendResponseHeaders(status.code, payloadSize)

        withContext(Dispatchers.IO) {
            IOUtils.copy(payload, exchange.responseBody)
        }
    }

    override suspend fun headers(): Map<String, List<String>> =
        exchange.requestHeaders

    override suspend fun queryParams(): Map<String, String> {
        val tokenizer = StringTokenizer(exchange.requestURI.query, "&")
        val output = mutableMapOf<String, String>()

        withContext(Dispatchers.IO) {
            while (tokenizer.hasMoreTokens()) {
                val token = tokenizer.nextToken()
                if (token.contains("=")) {
                    val (key, value) = token.split("=")
                    output[key] = URLDecoder.decode(value, Charsets.UTF_8.name())
                }
            }
        }

        return output
    }

    override suspend fun pathVariables(): Map<String, String> =
        pathVariables

    override suspend fun principal(): Subject? =
        principal

    override suspend fun uri(): String =
        exchange.requestURI.path

    private suspend fun invokeInterceptors() {
        interceptors.forEach { it.intercept(this, exchange) }
    }
}
