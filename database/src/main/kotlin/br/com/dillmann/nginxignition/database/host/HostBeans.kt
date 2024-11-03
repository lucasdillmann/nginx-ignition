package br.com.dillmann.nginxignition.database.host

import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxignition.core.host.HostRepository
import br.com.dillmann.nginxignition.database.host.lifecycle.HostMigrations
import org.koin.core.module.Module
import org.koin.dsl.bind

internal fun Module.hostBeans() {
    single { HostMigrations() } bind StartupCommand::class
    single { HostDatabaseRepository(get()) } bind HostRepository::class
    single { HostConverter() }
}
