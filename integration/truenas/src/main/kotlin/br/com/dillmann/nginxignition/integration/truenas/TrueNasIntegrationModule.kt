package br.com.dillmann.nginxignition.integration.truenas

import br.com.dillmann.nginxignition.core.integration.IntegrationAdapter
import org.koin.dsl.bind
import org.koin.dsl.module

object TrueNasIntegrationModule {
    fun initialize() =
        module {
            single { TrueNasIntegrationAdapter() } bind IntegrationAdapter::class
        }
}
