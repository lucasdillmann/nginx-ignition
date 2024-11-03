package br.com.dillmann.nginxignition.application.controller.nginx.handler

import br.com.dillmann.nginxignition.application.controller.nginx.model.NginxConverter
import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.nginx.command.StartNginxCommand
import br.com.dillmann.nginxignition.core.nginx.exception.NginxCommandException
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

class NginxStartHandler(
    private val startCommand: StartNginxCommand,
    private val converter: NginxConverter,
) {
    private val logger = logger<NginxStartHandler>()

    suspend fun handle(call: RoutingCall) {
        try {
            startCommand.start()
            call.respond(HttpStatusCode.NoContent)
        } catch (ex: NginxCommandException) {
            logger.warn("Start command failed to complete", ex)

            val payload = converter.toResponse(ex)
            call.respond(HttpStatusCode.FailedDependency, payload)
        }
    }
}
