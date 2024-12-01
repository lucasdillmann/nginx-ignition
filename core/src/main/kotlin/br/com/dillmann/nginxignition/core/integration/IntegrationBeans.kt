package br.com.dillmann.nginxignition.core.integration

import br.com.dillmann.nginxignition.core.integration.command.ConfigureIntegrationByIdCommand
import br.com.dillmann.nginxignition.core.integration.command.GetIntegrationByIdCommand
import br.com.dillmann.nginxignition.core.integration.command.ListIntegrationOptionsCommand
import br.com.dillmann.nginxignition.core.integration.command.ListIntegrationsCommand
import org.koin.core.module.Module
import org.koin.dsl.binds

fun Module.integrationBeans() {
    single { IntegrationService(get(), getAll(), get())} binds arrayOf(
        ConfigureIntegrationByIdCommand::class,
        GetIntegrationByIdCommand::class,
        ListIntegrationOptionsCommand::class,
        ListIntegrationsCommand::class,
    )
    single { IntegrationValidator() }
}
