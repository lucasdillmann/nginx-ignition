package br.com.dillmann.nginxignition.database.user

import br.com.dillmann.nginxignition.core.user.UserRepository
import org.koin.core.module.Module
import org.koin.dsl.bind

internal fun Module.userBeans() {
    single { UserDatabaseRepository(get()) } bind UserRepository::class
    single { UserConverter() }
}
