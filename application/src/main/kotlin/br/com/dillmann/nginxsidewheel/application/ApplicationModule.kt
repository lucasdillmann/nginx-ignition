package br.com.dillmann.nginxsidewheel.application

import br.com.dillmann.nginxsidewheel.application.controller.host.handler.*
import br.com.dillmann.nginxsidewheel.application.controller.host.model.HostConverter
import org.koin.dsl.module
import org.mapstruct.factory.Mappers

object ApplicationModule {
    fun initialize() =
        module {
            single { Mappers.getMapper(HostConverter::class.java) }
            single { DeleteHostByIdHandler(get()) }
            single { GetHostByIdHandler(get(), get()) }
            single { ListHostsHandler(get(), get()) }
            single { PutHostByIdHandler(get(), get()) }
            single { PostHostHandler(get(), get()) }
        }
}
