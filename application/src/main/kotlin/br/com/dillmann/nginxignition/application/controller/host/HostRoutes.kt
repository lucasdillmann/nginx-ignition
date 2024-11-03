package br.com.dillmann.nginxignition.application.controller.host

import br.com.dillmann.nginxignition.application.controller.host.handler.*
import io.ktor.server.application.*
import io.ktor.server.routing.*
import org.koin.ktor.ext.inject

fun Application.hostRoutes() {
    val listHandler by inject<ListHostsHandler>()
    val getByIdHandler by inject<GetHostByIdHandler>()
    val putByIdHandler by inject<UpdateHostByIdHandler>()
    val deleteByIdHandler by inject<DeleteHostByIdHandler>()
    val postHandler by inject<CreateHostHandler>()

    routing {
        route("/api/hosts") {
            get { listHandler.handle(call) }
            get("/{id}") { getByIdHandler.handle(call) }
            put("/{id}") { putByIdHandler.handle(call) }
            delete("/{id}") { deleteByIdHandler.handle(call) }
            post { postHandler.handle(call) }
        }
    }
}
