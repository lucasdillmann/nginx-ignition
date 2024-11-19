package br.com.dillmann.nginxignition.api.host.handler

import br.com.dillmann.nginxignition.api.host.model.HostConverter
import br.com.dillmann.nginxignition.api.host.model.HostRequest
import br.com.dillmann.nginxignition.core.host.command.SaveHostCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.IdAwareRequestHandler
import br.com.dillmann.nginxignition.api.common.request.payload
import java.util.UUID

internal class UpdateHostByIdHandler(
    private val saveCommand: SaveHostCommand,
    private val converter: HostConverter,
): IdAwareRequestHandler {
    override suspend fun handle(call: ApiCall, id: UUID) {
        val payload = call.payload<HostRequest>()
        val host = converter.toDomainModel(payload).copy(id = id)
        saveCommand.save(host)
        call.respond(HttpStatus.NO_CONTENT)
    }
}
