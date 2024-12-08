package br.com.dillmann.nginxignition.api.settings.handler

import br.com.dillmann.nginxignition.api.common.request.ApiCall
import br.com.dillmann.nginxignition.api.common.request.HttpStatus
import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.api.common.request.payload
import br.com.dillmann.nginxignition.api.settings.SettingsConverter
import br.com.dillmann.nginxignition.api.settings.model.SettingsDto
import br.com.dillmann.nginxignition.core.settings.command.SaveSettingsCommand

internal class PutSettingsHandler(
    private val saveCommand: SaveSettingsCommand,
    private val converter: SettingsConverter,
) : RequestHandler {
    override suspend fun handle(call: ApiCall) {
        val payload = call.payload<SettingsDto>().let(converter::toDomain)
        saveCommand.save(payload)
        call.respond(HttpStatus.NO_CONTENT)
    }
}
