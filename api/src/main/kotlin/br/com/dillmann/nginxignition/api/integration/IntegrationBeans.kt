package br.com.dillmann.nginxignition.api.integration

import br.com.dillmann.nginxignition.api.common.routing.RouteProvider
import br.com.dillmann.nginxignition.api.integration.handler.GetIntegrationConfigurationHandler
import br.com.dillmann.nginxignition.api.integration.handler.GetIntegrationOptionsHandler
import br.com.dillmann.nginxignition.api.integration.handler.ListIntegrationsHandler
import br.com.dillmann.nginxignition.api.integration.handler.PutIntegrationConfigurationHandler
import org.koin.core.module.Module
import org.koin.dsl.bind
import org.mapstruct.factory.Mappers

internal fun Module.integrationBeans() {
    single { Mappers.getMapper(IntegrationConverter::class.java) }
    single { GetIntegrationConfigurationHandler(get(), get()) }
    single { ListIntegrationsHandler(get(), get()) }
    single { PutIntegrationConfigurationHandler(get()) }
    single { GetIntegrationOptionsHandler(get(), get()) }
    single { IntegrationRoutes(get(), get(), get(), get()) } bind RouteProvider::class
}
