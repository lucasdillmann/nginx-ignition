package br.com.dillmann.nginxignition.api.common.request.handler

import br.com.dillmann.nginxignition.api.common.request.ApiCall

fun interface RequestHandler {
    suspend fun handle(call: ApiCall)
}
