package br.com.dillmann.nginxignition.database.settings

import br.com.dillmann.nginxignition.core.settings.SettingsRepository
import org.koin.core.module.Module
import org.koin.dsl.bind

internal fun Module.settingsBeans() {
    single { SettingsConverter() }
    single { SettingsDatabaseRepository(get()) } bind SettingsRepository::class
}
