package br.com.dillmann.nginxignition.application.common.routing

import io.ktor.server.routing.*
import io.ktor.utils.io.*

@KtorDsl
fun Route.get(pattern: String, handler: RequestHandler) =
    get(pattern) { handler.handle(call) }

@KtorDsl
fun Route.post(pattern: String, handler: RequestHandler) =
    post(pattern) { handler.handle(call) }

@KtorDsl
fun Route.put(pattern: String, handler: RequestHandler) =
    put(pattern) { handler.handle(call) }

@KtorDsl
fun Route.delete(pattern: String, handler: RequestHandler) =
    delete(pattern) { handler.handle(call) }
