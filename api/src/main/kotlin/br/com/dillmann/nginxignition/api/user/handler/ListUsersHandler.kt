package br.com.dillmann.nginxignition.api.user.handler

import br.com.dillmann.nginxignition.api.user.UserConverter
import br.com.dillmann.nginxignition.core.user.command.ListUserCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.PageAwareRequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond

internal class ListUsersHandler(
    private val listCommand: ListUserCommand,
    private val converter: UserConverter,
): PageAwareRequestHandler {
    override suspend fun handle(call: ApiCall, pageNumber: Int, pageSize: Int, searchTerms: String?) {
        val page = listCommand.list(pageSize, pageNumber, searchTerms)
        val payload = converter.toResponse(page)
        call.respond(HttpStatus.OK, payload)
    }
}
