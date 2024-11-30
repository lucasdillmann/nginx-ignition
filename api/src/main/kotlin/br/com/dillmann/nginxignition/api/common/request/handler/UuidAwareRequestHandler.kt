package br.com.dillmann.nginxignition.api.common.request.handler

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import java.util.UUID

internal interface UuidAwareRequestHandler: IdAwareRequestHandler {
    override suspend fun handle(call: ApiCall, id: String) {
        val uuid = runCatching { id.let(UUID::fromString) }.getOrNull()
        if (uuid == null) {
            call.respond(HttpStatus.BAD_REQUEST)
            return
        }

        handle(call, uuid)
    }

    suspend fun handle(call: ApiCall, id: UUID)
}
