package br.com.dillmann.nginxignition.api.host.handler

import br.com.dillmann.nginxignition.api.host.HostConverter
import br.com.dillmann.nginxignition.api.host.model.HostRequest
import br.com.dillmann.nginxignition.core.host.command.SaveHostCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.api.common.request.payload
import java.util.UUID

internal class CreateHostHandler(
    private val saveCommand: SaveHostCommand,
    private val converter: HostConverter,
): RequestHandler {
    override suspend fun handle(call: ApiCall) {
        val payload = call.payload<HostRequest>()
        val host = converter.toDomainModel(payload).copy(id = UUID.randomUUID())
        saveCommand.save(host)
        call.respond(HttpStatus.NO_CONTENT)
    }
}
