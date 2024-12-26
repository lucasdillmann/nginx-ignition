package br.com.dillmann.nginxignition.application.router

import br.com.dillmann.nginxignition.application.router.adapter.NettyApiCallAdapter
import br.com.dillmann.nginxignition.application.router.exception.ExceptionHandler
import br.com.dillmann.nginxignition.core.common.log.logger
import io.netty.channel.ChannelFutureListener
import io.netty.channel.ChannelHandlerContext
import io.netty.handler.codec.http.DefaultFullHttpResponse
import io.netty.handler.codec.http.FullHttpRequest
import io.netty.handler.codec.http.HttpResponseStatus
import io.netty.handler.codec.http.HttpVersion

internal class RequestRouter(
    compiler: RequestRouteCompiler,
    private val interceptors: List<ResponseInterceptor>,
    private val exceptionHandler: ExceptionHandler,
) {
    private companion object {
        private val LOGGER = logger<RequestRouter>()
    }

    private val routes = compiler.compile()

    suspend fun route(context: ChannelHandlerContext, request: FullHttpRequest) {
        try {
            val basePath = request.uri().split("?").first()
            for ((method, pattern, handler) in routes) {
                if (method != request.method()) continue

                val matcher = pattern.matcher(basePath)
                if (!matcher.find()) continue

                val pathVariables = matcher
                    .namedGroups()
                    .keys
                    .associateWith{ matcher.group(it) }

                val call = NettyApiCallAdapter(context, request, interceptors, null, pathVariables)
                handler.handle(call)
                return
            }

            sendResponse(context, HttpResponseStatus.NOT_FOUND, request.protocolVersion())
        } catch (ex: Throwable) {
            val call = NettyApiCallAdapter(context, request, interceptors)
            exceptionHandler.handle(call, ex)
        }
    }

    fun route(context: ChannelHandlerContext, exception: Throwable) {
        LOGGER.warn("Request failed with an exception", exception)
        sendResponse(context, HttpResponseStatus.INTERNAL_SERVER_ERROR)
    }

    private fun sendResponse(
        context: ChannelHandlerContext,
        status: HttpResponseStatus,
        version: HttpVersion = HttpVersion.HTTP_1_1,
    ) {
        val response = DefaultFullHttpResponse(version, status)
        context.write(response).addListener(ChannelFutureListener.CLOSE)
    }
}
