package br.com.dillmann.nginxignition.api.accesslist.handler

import br.com.dillmann.nginxignition.api.accesslist.AccessListConverter
import br.com.dillmann.nginxignition.api.accesslist.model.AccessListRequest
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.api.common.request.payload
import br.com.dillmann.nginxignition.core.accesslist.command.SaveAccessListByCommand

internal class PostAccessListHandler(
    private val saveCommand: SaveAccessListByCommand,
    private val converter: AccessListConverter,
) : RequestHandler {
    override suspend fun handle(call: ApiCall) {
        val payload = call.payload<AccessListRequest>().let(converter::toDomain)
        saveCommand.save(payload)
        call.respond(HttpStatus.NO_CONTENT)
    }
}
