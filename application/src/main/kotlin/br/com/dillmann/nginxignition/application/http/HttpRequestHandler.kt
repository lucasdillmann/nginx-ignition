package br.com.dillmann.nginxignition.application.http

import br.com.dillmann.nginxignition.application.router.RequestRouter
import com.sun.net.httpserver.HttpExchange
import com.sun.net.httpserver.HttpHandler
import kotlinx.coroutines.runBlocking

internal class HttpRequestHandler(private val router: RequestRouter): HttpHandler {
    override fun handle(exchange: HttpExchange) {
        exchange.use {
            runBlocking {
                router.route(exchange)
            }
        }
    }
}
