package br.com.dillmann.nginxignition.application.configuration

import br.com.dillmann.nginxignition.core.common.configuration.ConfigurationProvider
import org.koin.core.module.Module
import org.koin.dsl.bind

fun Module.configurationBeans() {
    single { RootConfigurationProvider() } bind ConfigurationProvider::class
}
