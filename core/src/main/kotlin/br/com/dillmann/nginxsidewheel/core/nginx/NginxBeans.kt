package br.com.dillmann.nginxsidewheel.core.nginx

import br.com.dillmann.nginxsidewheel.core.common.startup.StartupCommand
import br.com.dillmann.nginxsidewheel.core.nginx.command.GetStatusNginxCommand
import br.com.dillmann.nginxsidewheel.core.nginx.command.ReloadNginxCommand
import br.com.dillmann.nginxsidewheel.core.nginx.command.StartNginxCommand
import br.com.dillmann.nginxsidewheel.core.nginx.command.StopNginxCommand
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
    single { NginxConfigurationFiles(get(), get()) }
    single { NginxProcessManager(get()) }
    single { NginxSemaphore() }
}
