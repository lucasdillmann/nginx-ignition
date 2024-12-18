package br.com.dillmann.nginxignition.api.accesslist

import br.com.dillmann.nginxignition.api.accesslist.handler.*
import br.com.dillmann.nginxignition.api.accesslist.handler.DeleteAccessListByIdHandler
import br.com.dillmann.nginxignition.api.accesslist.handler.GetAccessListByIdHandler
import br.com.dillmann.nginxignition.api.accesslist.handler.ListAccessListHandler
import br.com.dillmann.nginxignition.api.accesslist.handler.PostAccessListHandler
import br.com.dillmann.nginxignition.api.common.routing.RouteProvider
import org.koin.core.module.Module
import org.koin.dsl.bind
import org.mapstruct.factory.Mappers

internal fun Module.accessListBeans() {
    single { Mappers.getMapper(AccessListConverter::class.java) }
    single { DeleteAccessListByIdHandler(get()) }
    single { GetAccessListByIdHandler(get(), get()) }
    single { ListAccessListHandler(get(), get()) }
    single { PostAccessListHandler(get(), get()) }
    single { PutAccessListHandler(get(), get()) }
    single { AccessListRoutes(get(), get(), get(), get(), get()) } bind RouteProvider::class
}
