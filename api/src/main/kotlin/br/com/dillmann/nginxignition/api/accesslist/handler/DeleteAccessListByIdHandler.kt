package br.com.dillmann.nginxignition.api.accesslist.handler

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.UuidAwareRequestHandler
import br.com.dillmann.nginxignition.core.accesslist.command.DeleteAccessListByIdCommand
import java.util.UUID

internal class DeleteAccessListByIdHandler(
    private val deleteCommand: DeleteAccessListByIdCommand,
) : UuidAwareRequestHandler {
    override suspend fun handle(call: ApiCall, id: UUID) {
        deleteCommand.deleteById(id)
        call.respond(HttpStatus.NO_CONTENT)
    }
}
