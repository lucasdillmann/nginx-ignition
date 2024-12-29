package br.com.dillmann.nginxignition.application

import br.com.dillmann.nginxignition.application.configuration.configurationBeans
import br.com.dillmann.nginxignition.application.frontend.frontendBeans
import br.com.dillmann.nginxignition.application.http.httpBeans
import br.com.dillmann.nginxignition.application.rbac.rbacBeans
import br.com.dillmann.nginxignition.application.router.routerBeans
import org.koin.dsl.module

object ApplicationModule {
    fun initialize() =
        module {
            configurationBeans()
            frontendBeans()
            httpBeans()
            rbacBeans()
            routerBeans()
        }
}
