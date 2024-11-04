package br.com.dillmann.nginxignition.application.controller.user.handler

import br.com.dillmann.nginxignition.application.controller.user.model.UserConverter
import br.com.dillmann.nginxignition.core.user.command.ListUserCommand
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

class ListUsersHandler(
    private val listCommand: ListUserCommand,
    private val converter: UserConverter,
) {
    suspend fun handle(call: RoutingCall) {
        val pageSize = runCatching { call.request.queryParameters["pageSize"]?.toInt() }.getOrNull() ?: 10
        val pageNumber = runCatching { call.request.queryParameters["pageNumber"]?.toInt() }.getOrNull() ?: 0

        val page = listCommand.list(pageSize, pageNumber)
        val payload = converter.toResponse(page)
        call.respond(HttpStatusCode.OK, payload)
    }
}
