package br.com.dillmann.nginxignition.core.settings

import br.com.dillmann.nginxignition.core.settings.command.GetSettingsCommand
import br.com.dillmann.nginxignition.core.settings.command.SaveSettingsCommand
import org.koin.core.module.Module
import org.koin.dsl.binds

internal fun Module.settingsBeans() {
    single { SettingsService(get(), get(), get()) } binds arrayOf(
        GetSettingsCommand::class,
        SaveSettingsCommand::class,
    )
    single { SettingsValidator(get()) }
}
