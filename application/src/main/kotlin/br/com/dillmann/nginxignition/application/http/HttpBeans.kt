package br.com.dillmann.nginxignition.application.http

import br.com.dillmann.nginxignition.application.http.lifecycle.HttpContainerShutdown
import br.com.dillmann.nginxignition.application.http.lifecycle.HttpContainerStartup
import br.com.dillmann.nginxignition.core.common.lifecycle.ShutdownCommand
import br.com.dillmann.nginxignition.core.common.lifecycle.StartupCommand
import org.koin.core.module.Module
import org.koin.dsl.bind

fun Module.httpBeans() {
    single { HttpContainer(get(), get()) }
    single { HttpRequestHandler(get()) }
    single { HttpContainerStartup(get()) } bind StartupCommand::class
    single { HttpContainerShutdown(get()) } bind ShutdownCommand::class
}
