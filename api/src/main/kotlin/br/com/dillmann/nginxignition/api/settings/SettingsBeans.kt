package br.com.dillmann.nginxignition.api.settings

import br.com.dillmann.nginxignition.api.common.routing.RouteProvider
import br.com.dillmann.nginxignition.api.settings.handler.GetSettingsHandler
import br.com.dillmann.nginxignition.api.settings.handler.PutSettingsHandler
import org.koin.core.module.Module
import org.koin.dsl.bind
import org.mapstruct.factory.Mappers

internal fun Module.settingsBeans() {
    single { Mappers.getMapper(SettingsConverter::class.java) }
    single { GetSettingsHandler(get(), get()) }
    single { PutSettingsHandler(get(), get()) }
    single { SettingsRoutes(get(), get()) } bind RouteProvider::class
}
