package br.com.dillmann.nginxignition.application.controller.user.handler

import br.com.dillmann.nginxignition.application.common.routing.template.PageAwareRequestHandler
import br.com.dillmann.nginxignition.application.controller.user.model.UserConverter
import br.com.dillmann.nginxignition.core.user.command.ListUserCommand
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

class ListUsersHandler(
    private val listCommand: ListUserCommand,
    private val converter: UserConverter,
): PageAwareRequestHandler {
    override suspend fun handle(call: RoutingCall, pageNumber: Int, pageSize: Int) {
        val page = listCommand.list(pageSize, pageNumber)
        val payload = converter.toResponse(page)
        call.respond(HttpStatusCode.OK, payload)
    }
}
