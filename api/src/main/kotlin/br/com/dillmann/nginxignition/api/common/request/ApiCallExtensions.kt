package br.com.dillmann.nginxignition.api.common.request

import kotlin.reflect.typeOf

suspend inline fun <reified T: Any> ApiCall.payload(): T = payload(T::class)

suspend inline fun <reified T: Any> ApiCall.respond(status: HttpStatus, payload: T) =
    respond(status, payload, payload::class, typeOf<T>())
