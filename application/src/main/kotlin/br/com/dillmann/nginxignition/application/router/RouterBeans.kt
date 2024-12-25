package br.com.dillmann.nginxignition.application.router

import org.koin.core.module.Module

fun Module.routerBeans() {
    single { RequestRouteCompiler(getAll()) }
    single { RequestRouter(get(), getAll()) }
}
