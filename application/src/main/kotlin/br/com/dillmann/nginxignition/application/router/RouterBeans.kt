package br.com.dillmann.nginxignition.application.router

import br.com.dillmann.nginxignition.application.router.exception.ExceptionHandler
import org.koin.core.module.Module

internal fun Module.routerBeans() {
    single { RequestRouteCompiler(getAll()) }
    single { RequestRouter(get(), getAll(), get()) }
    single { ExceptionHandler() }
}
