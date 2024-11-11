package br.com.dillmann.nginxignition.core.nginx

import br.com.dillmann.nginxignition.core.common.lifecycle.ShutdownCommand
import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import br.com.dillmann.nginxignition.core.nginx.command.GetStatusNginxCommand
import br.com.dillmann.nginxignition.core.nginx.command.ReloadNginxCommand
import br.com.dillmann.nginxignition.core.nginx.command.StartNginxCommand
import br.com.dillmann.nginxignition.core.nginx.command.StopNginxCommand
import br.com.dillmann.nginxignition.core.nginx.configuration.NginxConfigurationFacade
import br.com.dillmann.nginxignition.core.nginx.configuration.NginxConfigurationFileProvider
import br.com.dillmann.nginxignition.core.nginx.configuration.provider.*
import br.com.dillmann.nginxignition.core.nginx.configuration.provider.MainConfigurationFileProvider
import br.com.dillmann.nginxignition.core.nginx.lifecycle.NginxShutdown
import br.com.dillmann.nginxignition.core.nginx.lifecycle.NginxStartup
import org.koin.core.module.Module
import org.koin.dsl.bind
import org.koin.dsl.binds

internal fun Module.nginxBeans() {
    single { NginxService(get(), get(), get()) } binds arrayOf(
        ReloadNginxCommand::class,
        StartNginxCommand::class,
        StopNginxCommand::class,
        GetStatusNginxCommand::class,
    )
    single { NginxStartup(get()) } bind StartupCommand::class
    single { NginxShutdown(get()) } bind ShutdownCommand::class
    single { NginxProcessManager(get()) }
    single { NginxSemaphore() }
    single { NginxConfigurationFacade(get(), getAll(), get()) }
    single { MainConfigurationFileProvider() } bind NginxConfigurationFileProvider::class
    single { MimeTypesConfigurationFileProvider() } bind NginxConfigurationFileProvider::class
    single { HostConfigurationFileProvider() } bind NginxConfigurationFileProvider::class
    single { HostCertificateFileProvider(get()) } bind NginxConfigurationFileProvider::class
}
