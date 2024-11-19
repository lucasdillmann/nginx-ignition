package br.com.dillmann.nginxignition.api.host.handler

import br.com.dillmann.nginxignition.api.host.model.HostConverter
import br.com.dillmann.nginxignition.core.host.command.ListHostCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.PageAwareRequestHandler

internal class ListHostsHandler(
    private val listCommand: ListHostCommand,
    private val converter: HostConverter,
): PageAwareRequestHandler {
    override suspend fun handle(call: ApiCall, pageNumber: Int, pageSize: Int) {
        val page = listCommand.list(pageSize, pageNumber)
        val payload = converter.toResponse(page)
        call.respond(HttpStatus.OK, payload)
    }
}
