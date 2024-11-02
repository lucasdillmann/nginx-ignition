package br.com.dillmann.nginxsidewheel.database

import br.com.dillmann.nginxsidewheel.core.common.startup.StartupCommand
import br.com.dillmann.nginxsidewheel.database.common.DatabaseStartup
import br.com.dillmann.nginxsidewheel.database.host.hostBeans
import org.koin.dsl.bind
import org.koin.dsl.module

object DatabaseModule {
    fun initialize() =
        module {
            single { DatabaseStartup(get()) } bind StartupCommand::class
            hostBeans()
        }
}
