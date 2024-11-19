package br.com.dillmann.nginxignition.api.common.request

suspend inline fun <reified T: Any> ApiCall.payload(): T = payload(T::class)
