package br.com.dillmann.nginxignition.core.integration

import br.com.dillmann.nginxignition.core.integration.command.*
import org.koin.core.module.Module
import org.koin.dsl.binds

fun Module.integrationBeans() {
    single { IntegrationService(get(), getAll(), get())} binds arrayOf(
        ConfigureIntegrationByIdCommand::class,
        GetIntegrationByIdCommand::class,
        GetIntegrationOptionByIdCommand::class,
        ListIntegrationOptionsCommand::class,
        ListIntegrationsCommand::class,
    )
    single { IntegrationValidator() }
}
