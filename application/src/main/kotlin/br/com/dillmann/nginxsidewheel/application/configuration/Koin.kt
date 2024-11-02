package br.com.dillmann.nginxsidewheel.application.configuration

import br.com.dillmann.nginxsidewheel.application.ApplicationModule
import br.com.dillmann.nginxsidewheel.core.CoreModule
import br.com.dillmann.nginxsidewheel.database.DatabaseModule
import br.com.dillmann.nginxsidewheel.thirdparty.ThirdPartyModule
import io.ktor.server.application.*
import org.koin.ktor.plugin.Koin
import org.koin.logger.slf4jLogger

fun Application.configureKoin() {
    install(Koin) {
        slf4jLogger()
        modules(
            CoreModule.initialize(),
            DatabaseModule.initialize(),
            ThirdPartyModule.initialize(),
            ApplicationModule.initialize(),
        )
    }
}
