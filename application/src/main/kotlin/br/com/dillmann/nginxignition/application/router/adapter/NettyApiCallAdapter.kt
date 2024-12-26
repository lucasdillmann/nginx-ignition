package br.com.dillmann.nginxignition.application.router.adapter

import br.com.dillmann.nginxignition.api.common.authorization.Subject
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.application.router.ResponseInterceptor
import io.netty.buffer.Unpooled
import io.netty.channel.ChannelFutureListener
import io.netty.channel.ChannelHandlerContext
import io.netty.handler.codec.http.*
import kotlinx.serialization.json.Json
import kotlinx.serialization.serializer
import kotlin.reflect.KClass
import kotlin.reflect.KType
import io.netty.handler.codec.http.HttpHeaderNames.CONNECTION
import io.netty.handler.codec.http.HttpHeaderNames.CONTENT_LENGTH
import io.netty.handler.codec.http.HttpHeaderNames.CONTENT_TYPE
import io.netty.handler.codec.http.HttpHeaderValues.*

internal data class NettyApiCallAdapter(
    private val context: ChannelHandlerContext,
    private val request: FullHttpRequest,
    private val interceptors: List<ResponseInterceptor>,
    private val principal: Subject? = null,
    private val pathVariables: Map<String, String> = emptyMap(),

): ApiCall {
    @Suppress("UNCHECKED_CAST")
    override suspend fun <T : Any> payload(contract: KClass<T>): T =
        Json.decodeFromString(
            Json.serializersModule.serializer(contract.javaObjectType),
            request.content().toString(Charsets.UTF_8),
        ) as T

    override suspend fun <T : Any> respond(status: HttpStatus, payload: T, clazz: KClass<out T>, type: KType) {
        val json = Json.encodeToString(serializer(type), payload).encodeToByteArray()
        sendResponse(status, emptyMap(), json)
    }

    override suspend fun respond(status: HttpStatus) {
        sendResponse(status)
    }

    override suspend fun headers(): Map<String, List<String>> =
        request.headers().groupBy { it.key.lowercase() }.mapValues { pair -> pair.value.map { it.value } }

    override suspend fun queryParams(): Map<String, String> =
        QueryStringDecoder(request.uri()).parameters().mapValues { it.value.first() }

    override suspend fun pathVariables(): Map<String, String> =
        pathVariables

    override suspend fun principal(): Subject? =
        principal

    override suspend fun uri(): String =
        request.uri().split("?").first()

    override suspend fun respondRaw(status: HttpStatus, headers: Map<String, String>, payload: ByteArray?) {
        sendResponse(status, headers, payload)
    }

    private suspend fun sendResponse(
        status: HttpStatus,
        headers: Map<String, String> = emptyMap(),
        payload: ByteArray? = null,
    ) {
        val keepAlive = HttpUtil.isKeepAlive(request)
        val buffer = payload?.let(Unpooled::wrappedBuffer)
        val response = DefaultFullHttpResponse(
            request.protocolVersion(),
            HttpResponseStatus(status.code, status.name),
            buffer ?: Unpooled.EMPTY_BUFFER,
        )

        val responseHeaders = response.headers()
        if (headers.isNotEmpty()) {
            headers.forEach { responseHeaders.add(it.key, it.value) }
        } else {
            responseHeaders
                .add(CONTENT_TYPE, APPLICATION_JSON)
                .addInt(CONTENT_LENGTH, response.content().readableBytes())
        }

        if (keepAlive) {
            if (!request.protocolVersion().isKeepAliveDefault) {
                responseHeaders.add(CONNECTION, KEEP_ALIVE)
            }
        } else {
            responseHeaders.add(CONNECTION, CLOSE)
        }

        val interceptedResponse = invokeInterceptors(response)
        val future = context.write(interceptedResponse)
        if (!keepAlive) {
            future.addListener(ChannelFutureListener.CLOSE)
        }
    }

    private suspend fun invokeInterceptors(response: FullHttpResponse): FullHttpResponse {
        var output : FullHttpResponse = response
        interceptors.forEach { output = it.intercept(this, output) }
        return output
    }
}
