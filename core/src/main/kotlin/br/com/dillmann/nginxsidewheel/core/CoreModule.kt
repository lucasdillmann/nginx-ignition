package br.com.dillmann.nginxsidewheel.core

import br.com.dillmann.nginxsidewheel.core.host.hostBeans
import org.koin.dsl.module

object CoreModule {
    fun initialize() =
        module {
            hostBeans()
        }
}
