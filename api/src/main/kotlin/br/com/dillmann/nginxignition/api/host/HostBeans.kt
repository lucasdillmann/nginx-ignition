package br.com.dillmann.nginxignition.api.host

import br.com.dillmann.nginxignition.api.common.routing.RouteProvider
import br.com.dillmann.nginxignition.api.host.handler.*
import br.com.dillmann.nginxignition.api.host.model.HostConverter
import org.koin.core.module.Module
import org.koin.dsl.bind
import org.mapstruct.factory.Mappers

internal fun Module.hostBeans() {
    single { Mappers.getMapper(HostConverter::class.java) }
    single { DeleteHostByIdHandler(get()) }
    single { GetHostByIdHandler(get(), get()) }
    single { ListHostsHandler(get(), get()) }
    single { UpdateHostByIdHandler(get(), get()) }
    single { CreateHostHandler(get(), get()) }
    single { ToggleHostEnabledByIdHandler(get(), get()) }
    single { HostRoutes(get(), get(), get(), get(), get(), get()) } bind RouteProvider::class
}
