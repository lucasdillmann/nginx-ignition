package br.com.dillmann.nginxignition.application.controller.nginx.handler

import br.com.dillmann.nginxignition.application.controller.nginx.model.NginxStatusResponse
import br.com.dillmann.nginxignition.core.nginx.command.GetStatusNginxCommand
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

class NginxStatusHandler(
    private val getStatusCommand: GetStatusNginxCommand,
) {
    suspend fun handle(call: RoutingCall) {
        val running = getStatusCommand.isRunning()
        val payload = NginxStatusResponse(running)
        call.respond(HttpStatusCode.OK, payload)
    }
}
