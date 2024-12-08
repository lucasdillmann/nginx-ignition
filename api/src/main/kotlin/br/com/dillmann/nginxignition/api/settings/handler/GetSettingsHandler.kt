package br.com.dillmann.nginxignition.api.settings.handler

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.api.common.request.respond
import br.com.dillmann.nginxignition.api.settings.SettingsConverter
import br.com.dillmann.nginxignition.core.settings.command.GetSettingsCommand

internal class GetSettingsHandler(
    private val getCommand: GetSettingsCommand,
    private val converter: SettingsConverter,
) : RequestHandler {
    override suspend fun handle(call: ApiCall) {
        val payload = getCommand.get().let(converter::toResponse)
        call.respond(HttpStatus.OK, payload)
    }
}
