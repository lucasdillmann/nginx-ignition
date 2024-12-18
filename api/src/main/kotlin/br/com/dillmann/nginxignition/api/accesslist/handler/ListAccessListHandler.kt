package br.com.dillmann.nginxignition.api.accesslist.handler

import br.com.dillmann.nginxignition.api.accesslist.AccessListConverter
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.PageAwareRequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond
import br.com.dillmann.nginxignition.core.accesslist.command.ListAccessListCommand

internal class ListAccessListHandler(
    private val listCommand: ListAccessListCommand,
    private val converter: AccessListConverter,
) : PageAwareRequestHandler {
    override suspend fun handle(call: ApiCall, pageNumber: Int, pageSize: Int, searchTerms: String?) {
        val payload = listCommand.getPage(pageSize, pageNumber, searchTerms).let(converter::toResponse)
        call.respond(HttpStatus.OK, payload)
    }
}
