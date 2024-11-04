package br.com.dillmann.nginxignition.application.controller.nginx

import br.com.dillmann.nginxignition.application.common.rbac.requireAuthentication
import br.com.dillmann.nginxignition.application.controller.nginx.handler.NginxReloadHandler
import br.com.dillmann.nginxignition.application.controller.nginx.handler.NginxStartHandler
import br.com.dillmann.nginxignition.application.controller.nginx.handler.NginxStatusHandler
import br.com.dillmann.nginxignition.application.controller.nginx.handler.NginxStopHandler
import io.ktor.server.application.*
import io.ktor.server.routing.*
import org.koin.ktor.ext.inject

fun Application.nginxRoutes() {
    val startHandler by inject<NginxStartHandler>()
    val stopHandler by inject<NginxStopHandler>()
    val reloadHandler by inject<NginxReloadHandler>()
    val statusHandler by inject<NginxStatusHandler>()

    routing {
        requireAuthentication {
            route("/api/nginx") {
                post("/start") { startHandler.handle(call) }
                post("/stop") { stopHandler.handle(call) }
                post("/reload") { reloadHandler.handle(call) }
                get("/status") { statusHandler.handle(call) }
            }
        }
    }
}
