package br.com.dillmann.nginxignition.database.accesslist

import br.com.dillmann.nginxignition.core.accesslist.AccessListRepository
import org.koin.core.module.Module
import org.koin.dsl.bind

internal fun Module.accessListBeans() {
    single { AccessListDatabaseRepository(get()) } bind AccessListRepository::class
    single { AccessListConverter() }
}
