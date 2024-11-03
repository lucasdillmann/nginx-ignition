package br.com.dillmann.nginxignition.database.user

import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxignition.core.user.UserRepository
import br.com.dillmann.nginxignition.database.user.lifecycle.UserMigrations
import org.koin.core.module.Module
import org.koin.dsl.bind

internal fun Module.userBeans() {
    single { UserMigrations() } bind StartupCommand::class
    single { UserDatabaseRepository(get()) } bind UserRepository::class
    single { UserConverter() }
}
