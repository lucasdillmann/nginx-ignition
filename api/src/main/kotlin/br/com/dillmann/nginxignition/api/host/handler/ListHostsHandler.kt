package br.com.dillmann.nginxignition.api.host.handler

import br.com.dillmann.nginxignition.api.host.HostConverter
import br.com.dillmann.nginxignition.core.host.command.ListHostCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.PageAwareRequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond

internal class ListHostsHandler(
    private val listCommand: ListHostCommand,
    private val converter: HostConverter,
): PageAwareRequestHandler {
    override suspend fun handle(call: ApiCall, pageNumber: Int, pageSize: Int, searchTerms: String?) {
        val page = listCommand.list(pageSize, pageNumber, searchTerms)
        val payload = converter.toResponse(page)
        call.respond(HttpStatus.OK, payload)
    }
}
