package br.com.dillmann.nginxignition.database.host

import br.com.dillmann.nginxignition.core.host.HostRepository
import org.koin.core.module.Module
import org.koin.dsl.bind

internal fun Module.hostBeans() {
    single { HostDatabaseRepository(get()) } bind HostRepository::class
    single { HostConverter() }
}
