package br.com.dillmann.nginxignition.api.accesslist.handler

import br.com.dillmann.nginxignition.api.accesslist.AccessListConverter
import br.com.dillmann.nginxignition.api.accesslist.model.AccessListRequest
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.UuidAwareRequestHandler
import br.com.dillmann.nginxignition.api.common.request.payload
import br.com.dillmann.nginxignition.core.accesslist.command.SaveAccessListByCommand
import java.util.UUID

internal class PutAccessListHandler(
    private val saveCommand: SaveAccessListByCommand,
    private val converter: AccessListConverter,
) : UuidAwareRequestHandler {
    override suspend fun handle(call: ApiCall, id: UUID) {
        val payload = call.payload<AccessListRequest>().let(converter::toDomain).copy(id = id)
        saveCommand.save(payload)
        call.respond(HttpStatus.NO_CONTENT)
    }
}
