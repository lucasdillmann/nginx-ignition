package br.com.dillmann.nginxignition.api.common.request.handler

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus

internal interface IdAwareRequestHandler: RequestHandler {
    override suspend fun handle(call: ApiCall) {
        val id = runCatching { call.pathVariables()["id"] }.getOrNull()
        if (id == null) {
            call.respond(HttpStatus.BAD_REQUEST)
            return
        }

        handle(call, id)
    }

    suspend fun handle(call: ApiCall, id: String)
}
