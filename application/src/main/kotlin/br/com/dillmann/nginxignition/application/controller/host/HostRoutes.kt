package br.com.dillmann.nginxignition.application.controller.host

import br.com.dillmann.nginxignition.application.common.rbac.requireAuthentication
import br.com.dillmann.nginxignition.application.common.routing.*
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
        requireAuthentication {
            route("/api/hosts") {
                get(listHandler)
                get("/{id}", getByIdHandler)
                put("/{id}", putByIdHandler)
                delete("/{id}", deleteByIdHandler)
                post(postHandler)
                get("/{id}/access-logs") { TODO() }
                get("/{id}/error-logs") { TODO() }
            }
        }
    }
}
