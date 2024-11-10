package br.com.dillmann.nginxignition.application.controller.host.handler

import br.com.dillmann.nginxignition.application.common.routing.template.PageAwareRequestHandler
import br.com.dillmann.nginxignition.application.controller.host.model.HostConverter
import br.com.dillmann.nginxignition.core.host.command.ListHostCommand
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

class ListHostsHandler(
    private val listCommand: ListHostCommand,
    private val converter: HostConverter,
): PageAwareRequestHandler {
    override suspend fun handle(call: RoutingCall, pageNumber: Int, pageSize: Int) {
        val page = listCommand.list(pageSize, pageNumber)
        val payload = converter.toResponse(page)
        call.respond(HttpStatusCode.OK, payload)
    }
}
