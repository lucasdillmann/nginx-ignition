package br.com.dillmann.nginxignition.application

import br.com.dillmann.nginxignition.application.common.lifecycle.LifecycleManager
import br.com.dillmann.nginxignition.application.common.provider.CompositeConfigurationProvider
import br.com.dillmann.nginxignition.application.common.rbac.RbacJwtFacade
import br.com.dillmann.nginxignition.application.controller.host.hostBeans
import br.com.dillmann.nginxignition.application.controller.nginx.nginxBeans
import br.com.dillmann.nginxignition.application.controller.user.userBeans
import br.com.dillmann.nginxignition.core.common.provider.ConfigurationProvider
import org.koin.dsl.bind
import org.koin.dsl.module

object ApplicationModule {
    fun initialize() =
        module {
            single { CompositeConfigurationProvider() } bind ConfigurationProvider::class
            single { LifecycleManager(getAll(), getAll()) }
            single { RbacJwtFacade(get(), get(), get()) }

            hostBeans()
            nginxBeans()
            userBeans()
        }
}
