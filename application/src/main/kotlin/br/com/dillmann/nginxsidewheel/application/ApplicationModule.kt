package br.com.dillmann.nginxsidewheel.application

import br.com.dillmann.nginxsidewheel.application.common.lifecycle.LifecycleManager
import br.com.dillmann.nginxsidewheel.application.common.provider.CompositeConfigurationProvider
import br.com.dillmann.nginxsidewheel.application.controller.host.hostBeans
import br.com.dillmann.nginxsidewheel.application.controller.nginx.nginxBeans
import br.com.dillmann.nginxsidewheel.core.common.provider.ConfigurationProvider
import org.koin.dsl.bind
import org.koin.dsl.module

object ApplicationModule {
    fun initialize() =
        module {
            single { CompositeConfigurationProvider() } bind ConfigurationProvider::class
            single { LifecycleManager(getAll(), getAll()) }

            hostBeans()
            nginxBeans()
        }
}
