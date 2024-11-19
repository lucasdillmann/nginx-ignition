package br.com.dillmann.nginxignition.api.host.handler

import br.com.dillmann.nginxignition.core.host.command.DeleteHostCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.IdAwareRequestHandler
import java.util.UUID

internal class DeleteHostByIdHandler(
    private val deleteCommand: DeleteHostCommand,
): IdAwareRequestHandler {
    override suspend fun handle(call: ApiCall, id: UUID) {
        deleteCommand.deleteById(id)
        call.respond(HttpStatus.NO_CONTENT)
    }
}
