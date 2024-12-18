package br.com.dillmann.nginxignition.api.accesslist.handler

import br.com.dillmann.nginxignition.api.accesslist.AccessListConverter
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.UuidAwareRequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond
import br.com.dillmann.nginxignition.core.accesslist.command.GetAccessListByIdCommand
import java.util.UUID

internal class GetAccessListByIdHandler(
    private val getCommand: GetAccessListByIdCommand,
    private val converter: AccessListConverter,
) : UuidAwareRequestHandler {
    override suspend fun handle(call: ApiCall, id: UUID) {
        val payload = getCommand.getById(id)?.let(converter::toResponse)
        if (payload == null)
            call.respond(HttpStatus.NOT_FOUND)
        else
            call.respond(HttpStatus.OK, payload)
    }
}
