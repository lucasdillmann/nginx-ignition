package br.com.dillmann.nginxignition.application.router

import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.application.router.adapter.HttpApiCallAdapter
import br.com.dillmann.nginxignition.application.router.exception.ExceptionHandler
import com.sun.net.httpserver.HttpExchange

internal class RequestRouter(
    compiler: RequestRouteCompiler,
    private val interceptors: List<ResponseInterceptor>,
    private val exceptionHandler: ExceptionHandler,
) {
    private val routes = compiler.compile()

    suspend fun route(exchange: HttpExchange) {
        try {
            val basePath = exchange.requestURI.path
            for ((method, pattern, handler) in routes) {
                if (method != exchange.requestMethod) continue

                val matcher = pattern.matcher(basePath)
                if (!matcher.find()) continue

                val pathVariables = matcher
                    .namedGroups()
                    .keys
                    .associateWith{ matcher.group(it) }

                val call = HttpApiCallAdapter(exchange, interceptors, null, pathVariables)
                handler.handle(call)
                return
            }

            exchange.sendResponseHeaders(HttpStatus.NOT_FOUND.code, 0)
        } catch (ex: Throwable) {
            val call = HttpApiCallAdapter(exchange, interceptors)
            exceptionHandler.handle(call, ex)
        }
    }
}
