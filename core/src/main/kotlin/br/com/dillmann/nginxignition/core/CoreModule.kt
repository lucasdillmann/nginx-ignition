package br.com.dillmann.nginxignition.core

import br.com.dillmann.nginxignition.core.host.hostBeans
import br.com.dillmann.nginxignition.core.nginx.nginxBeans
import br.com.dillmann.nginxignition.core.user.userBeans
import org.koin.dsl.module

object CoreModule {
    fun initialize() =
        module {
            hostBeans()
            nginxBeans()
            userBeans()
        }
}
