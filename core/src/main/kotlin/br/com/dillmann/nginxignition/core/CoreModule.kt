package br.com.dillmann.nginxignition.core

import br.com.dillmann.nginxignition.core.accesslist.accessListBeans
import br.com.dillmann.nginxignition.core.certificate.certificateBeans
import br.com.dillmann.nginxignition.core.common.commonBeans
import br.com.dillmann.nginxignition.core.host.hostBeans
import br.com.dillmann.nginxignition.core.integration.integrationBeans
import br.com.dillmann.nginxignition.core.nginx.nginxBeans
import br.com.dillmann.nginxignition.core.settings.settingsBeans
import br.com.dillmann.nginxignition.core.user.userBeans
import org.koin.dsl.module

object CoreModule {
    fun initialize() =
        module {
            accessListBeans()
            commonBeans()
            certificateBeans()
            hostBeans()
            integrationBeans()
            nginxBeans()
            settingsBeans()
            userBeans()
        }
}
