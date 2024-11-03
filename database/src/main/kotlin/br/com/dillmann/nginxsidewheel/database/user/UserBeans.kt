package br.com.dillmann.nginxsidewheel.database.user

import br.com.dillmann.nginxsidewheel.core.common.startup.StartupCommand
import br.com.dillmann.nginxsidewheel.core.user.UserRepository
import org.koin.core.module.Module
import org.koin.dsl.bind

internal fun Module.userBeans() {
    single { UserMigrations() } bind StartupCommand::class
    single { UserDatabaseRepository(get()) } bind UserRepository::class
    single { UserConverter() }
}
