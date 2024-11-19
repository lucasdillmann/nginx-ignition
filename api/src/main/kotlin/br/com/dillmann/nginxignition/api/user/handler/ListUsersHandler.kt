package br.com.dillmann.nginxignition.api.user.handler

import br.com.dillmann.nginxignition.api.user.model.UserConverter
import br.com.dillmann.nginxignition.core.user.command.ListUserCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.PageAwareRequestHandler

internal class ListUsersHandler(
    private val listCommand: ListUserCommand,
    private val converter: UserConverter,
): PageAwareRequestHandler {
    override suspend fun handle(call: ApiCall, pageNumber: Int, pageSize: Int) {
        val page = listCommand.list(pageSize, pageNumber)
        val payload = converter.toResponse(page)
        call.respond(HttpStatus.OK, payload)
    }
}
