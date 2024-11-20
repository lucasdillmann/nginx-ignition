package br.com.dillmann.nginxignition.api.nginx.handler

import br.com.dillmann.nginxignition.api.nginx.model.NginxStatusResponse
import br.com.dillmann.nginxignition.core.nginx.command.GetStatusNginxCommand
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond

internal class NginxStatusHandler(
    private val getStatusCommand: GetStatusNginxCommand,
): RequestHandler {
    override suspend fun handle(call: ApiCall) {
        val running = getStatusCommand.isRunning()
        val payload = NginxStatusResponse(running)
        call.respond(HttpStatus.OK, payload)
    }
}
