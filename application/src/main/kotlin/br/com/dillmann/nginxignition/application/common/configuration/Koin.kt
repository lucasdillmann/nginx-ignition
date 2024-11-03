package br.com.dillmann.nginxignition.application.common.configuration

import br.com.dillmann.nginxignition.application.ApplicationModule
import br.com.dillmann.nginxignition.core.CoreModule
import br.com.dillmann.nginxignition.database.DatabaseModule
import br.com.dillmann.nginxignition.thirdparty.ThirdPartyModule
import io.ktor.server.application.*
import org.koin.ktor.plugin.Koin

fun Application.configureKoin() {
    install(Koin) {
        modules(
            CoreModule.initialize(),
            DatabaseModule.initialize(),
            ThirdPartyModule.initialize(),
            ApplicationModule.initialize(),
        )
    }
}
