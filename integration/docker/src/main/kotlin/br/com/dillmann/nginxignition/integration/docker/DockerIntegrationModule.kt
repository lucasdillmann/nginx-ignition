package br.com.dillmann.nginxignition.integration.docker

import br.com.dillmann.nginxignition.core.integration.IntegrationAdapter
import org.koin.dsl.bind
import org.koin.dsl.module

object DockerIntegrationModule {
    fun initialize() =
        module {
            single { DockerIntegrationAdapter() } bind IntegrationAdapter::class
        }
}
