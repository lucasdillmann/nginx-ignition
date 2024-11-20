package br.com.dillmann.nginxignition.api.nginx.handler

import br.com.dillmann.nginxignition.api.nginx.model.NginxConverter
import br.com.dillmann.nginxignition.core.common.log.logger
import br.com.dillmann.nginxignition.core.nginx.command.StartNginxCommand
import br.com.dillmann.nginxignition.core.nginx.exception.NginxCommandException
import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond

internal class NginxStartHandler(
    private val startCommand: StartNginxCommand,
    private val converter: NginxConverter,
): RequestHandler {
    private val logger = logger<NginxStartHandler>()

    override suspend fun handle(call: ApiCall) {
        try {
            startCommand.start()
            call.respond(HttpStatus.NO_CONTENT)
        } catch (ex: NginxCommandException) {
            logger.warn("Start command failed to complete", ex)

            val payload = converter.toResponse(ex)
            call.respond(HttpStatus.FAILED_DEPENDENCY, payload)
        }
    }
}
