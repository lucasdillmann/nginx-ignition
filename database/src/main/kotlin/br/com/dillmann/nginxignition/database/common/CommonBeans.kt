package br.com.dillmann.nginxignition.database.common

import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxignition.database.common.migrations.MigrationsService
import br.com.dillmann.nginxignition.database.common.migrations.MigrationsStartup
import org.koin.core.module.Module
import org.koin.dsl.bind

internal fun Module.commonBeans() {
    single { MigrationsStartup(get()) } bind StartupCommand::class
    single { MigrationsService() }
}
