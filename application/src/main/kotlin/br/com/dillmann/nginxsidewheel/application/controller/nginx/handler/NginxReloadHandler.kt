package br.com.dillmann.nginxsidewheel.application.controller.nginx.handler

import br.com.dillmann.nginxsidewheel.application.controller.nginx.model.NginxConverter
import br.com.dillmann.nginxsidewheel.core.common.log.logger
import br.com.dillmann.nginxsidewheel.core.nginx.command.ReloadNginxCommand
import br.com.dillmann.nginxsidewheel.core.nginx.exception.NginxCommandException
import io.ktor.http.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

class NginxReloadHandler(
    private val reloadCommand: ReloadNginxCommand,
    private val converter: NginxConverter,
) {
    private val logger = logger<NginxReloadHandler>()

    suspend fun handle(call: RoutingCall) {
        try {
            reloadCommand.reload()
            call.respond(HttpStatusCode.NoContent)
        } catch (ex: NginxCommandException) {
            logger.warn("Reload command failed to complete", ex)

            val payload = converter.toResponse(ex)
            call.respond(HttpStatusCode.FailedDependency, payload)
        }
    }
}
