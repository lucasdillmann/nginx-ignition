package br.com.dillmann.nginxsidewheel.database

import br.com.dillmann.nginxsidewheel.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxsidewheel.database.common.lifecycle.DatabaseStartup
import br.com.dillmann.nginxsidewheel.database.host.hostBeans
import br.com.dillmann.nginxsidewheel.database.user.userBeans
import org.koin.dsl.bind
import org.koin.dsl.module

object DatabaseModule {
    fun initialize() =
        module {
            single { DatabaseStartup(get()) } bind StartupCommand::class
            hostBeans()
            userBeans()
        }
}
