package br.com.dillmann.nginxignition.api.nginx

import br.com.dillmann.nginxignition.api.nginx.handler.NginxReloadHandler
import br.com.dillmann.nginxignition.api.nginx.handler.NginxStartHandler
import br.com.dillmann.nginxignition.api.nginx.handler.NginxStatusHandler
import br.com.dillmann.nginxignition.api.nginx.handler.NginxStopHandler
import br.com.dillmann.nginxignition.api.nginx.model.NginxConverter
import br.com.dillmann.nginxignition.api.common.routing.RouteProvider
import org.koin.core.module.Module
import org.koin.dsl.bind
import org.mapstruct.factory.Mappers

internal fun Module.nginxBeans() {
    single { Mappers.getMapper(NginxConverter::class.java) }
    single { NginxStartHandler(get(), get()) }
    single { NginxStopHandler(get(), get()) }
    single { NginxReloadHandler(get(), get()) }
    single { NginxStatusHandler(get()) }
    single { NginxRoutes(get(), get(), get(), get()) } bind RouteProvider::class
}
