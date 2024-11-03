package br.com.dillmann.nginxignition.application.controller.host.handler

import br.com.dillmann.nginxignition.application.controller.host.model.HostConverter
import br.com.dillmann.nginxignition.core.host.command.ListHostCommand
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

class ListHostsHandler(
    private val listCommand: ListHostCommand,
    private val converter: HostConverter,
) {
    suspend fun handle(call: RoutingCall) {
        val pageSize = runCatching { call.request.queryParameters["pageSize"]?.toInt() }.getOrNull() ?: 10
        val pageNumber = runCatching { call.request.queryParameters["pageNumber"]?.toInt() }.getOrNull() ?: 0

        val page = listCommand.list(pageSize, pageNumber)
        val payload = converter.toResponse(page)
        call.respond(HttpStatusCode.OK, payload)
    }
}
