package br.com.dillmann.nginxignition.api.host.handler

import br.com.dillmann.nginxignition.core.host.command.SaveHostCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.IdAwareRequestHandler
import br.com.dillmann.nginxignition.core.host.command.GetHostCommand
import java.util.UUID

internal class ToggleHostEnabledByIdHandler(
    private val getCommand: GetHostCommand,
    private val saveCommand: SaveHostCommand,
): IdAwareRequestHandler {
    override suspend fun handle(call: ApiCall, id: UUID) {
        val host = getCommand.getById(id)
        if (host == null) {
            call.respond(HttpStatus.NOT_FOUND)
            return
        }

        saveCommand.save(host.copy(enabled = !host.enabled))
        call.respond(HttpStatus.NO_CONTENT)
    }
}
