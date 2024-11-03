package br.com.dillmann.nginxsidewheel.application.controller.nginx

import br.com.dillmann.nginxsidewheel.application.controller.nginx.handler.*
import io.ktor.server.application.*
import io.ktor.server.routing.*
import org.koin.ktor.ext.inject

fun Application.nginxRoutes() {
    val startHandler by inject<NginxStartHandler>()
    val stopHandler by inject<NginxStopHandler>()
    val reloadHandler by inject<NginxReloadHandler>()
    val statusHandler by inject<NginxStatusHandler>()

    routing {
        route("/api/nginx") {
            post("/start") { startHandler.handle(call) }
            post("/stop") { stopHandler.handle(call) }
            post("/reload") { reloadHandler.handle(call) }
            get("/status") { statusHandler.handle(call) }
        }
    }
}
