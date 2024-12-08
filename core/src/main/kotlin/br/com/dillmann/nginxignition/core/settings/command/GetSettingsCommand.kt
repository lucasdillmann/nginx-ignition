package br.com.dillmann.nginxignition.core.settings.command

import br.com.dillmann.nginxignition.core.settings.Settings

fun interface GetSettingsCommand {
    suspend fun get(): Settings
}
