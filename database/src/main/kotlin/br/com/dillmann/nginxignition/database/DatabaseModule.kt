package br.com.dillmann.nginxignition.database

import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxignition.database.certificate.certificateBeans
import br.com.dillmann.nginxignition.database.common.lifecycle.DatabaseStartup
import br.com.dillmann.nginxignition.database.host.hostBeans
import br.com.dillmann.nginxignition.database.user.userBeans
import org.koin.dsl.bind
import org.koin.dsl.module

object DatabaseModule {
    fun initialize() =
        module {
            single { DatabaseStartup(get()) } bind StartupCommand::class

            certificateBeans()
            hostBeans()
            userBeans()
        }
}
