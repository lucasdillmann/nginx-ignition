package br.com.dillmann.nginxignition.application

import br.com.dillmann.nginxignition.application.exception.ConsistencyExceptionHandler
import br.com.dillmann.nginxignition.application.lifecycle.LifecycleManager
import br.com.dillmann.nginxignition.application.provider.CompositeConfigurationProvider
import br.com.dillmann.nginxignition.application.rbac.RbacJwtFacade
import br.com.dillmann.nginxignition.application.rbac.RbacJwtAuthorizer
import br.com.dillmann.nginxignition.core.common.configuration.ConfigurationProvider
import br.com.dillmann.nginxignition.api.common.authorization.Authorizer
import org.koin.dsl.bind
import org.koin.dsl.module

object ApplicationModule {
    fun initialize() =
        module {
            single { CompositeConfigurationProvider() } bind ConfigurationProvider::class
            single { RbacJwtAuthorizer(get(), get()) } bind Authorizer::class
            single { LifecycleManager(getAll(), getAll()) }
            single { RbacJwtFacade(get(), get(), get()) }
            single { ConsistencyExceptionHandler() }
        }
}
