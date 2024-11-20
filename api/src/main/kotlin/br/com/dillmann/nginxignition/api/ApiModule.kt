package br.com.dillmann.nginxignition.api

import br.com.dillmann.nginxignition.api.certificate.certificateBeans
import br.com.dillmann.nginxignition.api.host.hostBeans
import br.com.dillmann.nginxignition.api.nginx.nginxBeans
import br.com.dillmann.nginxignition.api.user.userBeans
import org.koin.dsl.module

object ApiModule {
    fun initialize() =
        module {
            certificateBeans()
            hostBeans()
            nginxBeans()
            userBeans()
        }
}