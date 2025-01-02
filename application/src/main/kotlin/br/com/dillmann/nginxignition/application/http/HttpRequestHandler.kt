package br.com.dillmann.nginxignition.application.http

import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.application.router.RequestRouter
import com.sun.net.httpserver.HttpExchange
import com.sun.net.httpserver.HttpHandler
import kotlinx.coroutines.runBlocking
import org.slf4j.LoggerFactory

internal class HttpRequestHandler(private val router: RequestRouter): HttpHandler {
    private companion object {
        private val LOGGER = LoggerFactory.getLogger(HttpRequestHandler::class.java)
    }

    override fun handle(exchange: HttpExchange) {
        exchange.use {
            try {
                runBlocking {
                    router.route(exchange)
                }
            } catch (ex: Exception) {
                LOGGER.error("Request failed with an unhandled exception", ex)
                it.sendResponseHeaders(HttpStatus.INTERNAL_SERVER_ERROR.code, -1)
            }
        }
    }
}
