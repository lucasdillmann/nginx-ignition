package br.com.dillmann.nginxignition.core.settings.command

import br.com.dillmann.nginxignition.core.settings.Settings

fun interface SaveSettingsCommand {
    suspend fun save(settings: Settings)
}
