package br.com.dillmann.nginxsidewheel.core

import br.com.dillmann.nginxsidewheel.core.host.hostBeans
import br.com.dillmann.nginxsidewheel.core.nginx.nginxBeans
import org.koin.dsl.module

object CoreModule {
    fun initialize() =
        module {
            hostBeans()
            nginxBeans()
        }
}
