package br.com.dillmann.nginxignition.api.accesslist.handler

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.UuidAwareRequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond
import br.com.dillmann.nginxignition.core.accesslist.command.DeleteAccessListByIdCommand
import java.util.UUID

internal class DeleteAccessListByIdHandler(
    private val deleteCommand: DeleteAccessListByIdCommand,
) : UuidAwareRequestHandler {
    override suspend fun handle(call: ApiCall, id: UUID) {
        val (deleted, reason) = deleteCommand.deleteById(id)
        if (!deleted) {
            val payload = mapOf("message" to reason)
            call.respond(HttpStatus.PRECONDITION_FAILED, payload)
            return
        }

        call.respond(HttpStatus.NO_CONTENT)
    }
}
