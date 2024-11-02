package br.com.dillmann.nginxsidewheel.application

import br.com.dillmann.nginxsidewheel.application.controller.host.hostBeans
import org.koin.dsl.module

object ApplicationModule {
    fun initialize() =
        module {
            hostBeans()
        }
}
