package br.com.dillmann.nginxignition.application.controller.host

import br.com.dillmann.nginxignition.application.controller.host.handler.*
import br.com.dillmann.nginxignition.application.controller.host.model.HostConverter
import org.koin.core.module.Module
import org.mapstruct.factory.Mappers

internal fun Module.hostBeans() {
    single { Mappers.getMapper(HostConverter::class.java) }
    single { DeleteHostByIdHandler(get()) }
    single { GetHostByIdHandler(get(), get()) }
    single { ListHostsHandler(get(), get()) }
    single { UpdateHostByIdHandler(get(), get()) }
    single { CreateHostHandler(get(), get()) }
}
