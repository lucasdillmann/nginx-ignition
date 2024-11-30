package br.com.dillmann.nginxignition.api.nginx.handler

import br.com.dillmann.nginxignition.api.nginx.NginxConverter
import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.nginx.command.StopNginxCommand
import br.com.dillmann.nginxignition.core.nginx.exception.NginxCommandException
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond

internal class NginxStopHandler(
    private val stopCommand: StopNginxCommand,
    private val converter: NginxConverter,
): RequestHandler {
    private val logger = logger<NginxStopHandler>()

    override suspend fun handle(call: ApiCall) {
        try {
            stopCommand.stop()
            call.respond(HttpStatus.NO_CONTENT)
        } catch (ex: NginxCommandException) {
            logger.warn("Stop command failed to complete", ex)

            val payload = converter.toResponse(ex)
            call.respond(HttpStatus.FAILED_DEPENDENCY, payload)
        }
    }
}
