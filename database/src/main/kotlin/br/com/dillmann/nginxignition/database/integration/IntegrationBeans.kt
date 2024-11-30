package br.com.dillmann.nginxignition.database.integration

import br.com.dillmann.nginxignition.core.integration.IntegrationRepository
import org.koin.core.module.Module
import org.koin.dsl.bind

internal fun Module.integrationBeans() {
    single { IntegrationDatabaseRepository(get()) } bind IntegrationRepository::class
    single { IntegrationConverter() }
}
