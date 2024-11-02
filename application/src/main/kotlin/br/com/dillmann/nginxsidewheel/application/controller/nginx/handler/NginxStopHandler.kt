package br.com.dillmann.nginxsidewheel.application.controller.host.handler

import br.com.dillmann.nginxsidewheel.application.controller.nginx.model.NginxConverter
import br.com.dillmann.nginxsidewheel.core.common.log.logger
import br.com.dillmann.nginxsidewheel.core.nginx.command.StopNginxCommand
import br.com.dillmann.nginxsidewheel.core.nginx.exception.NginxCommandException
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

class NginxStopHandler(
    private val stopCommand: StopNginxCommand,
    private val converter: NginxConverter,
) {
    private val logger = logger<NginxStopHandler>()

    suspend fun handle(call: RoutingCall) {
        try {
            stopCommand.stop()
            call.respond(HttpStatusCode.NoContent)
        } catch (ex: NginxCommandException) {
            logger.warn("Stop command failed to complete", ex)

            val payload = converter.toResponse(ex)
            call.respond(HttpStatusCode.FailedDependency, payload)
        }
    }
}
