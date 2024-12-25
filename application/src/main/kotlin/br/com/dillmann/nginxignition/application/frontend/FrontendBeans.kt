package br.com.dillmann.nginxignition.application.frontend

import br.com.dillmann.nginxignition.api.common.request.handler.RequestHandler
import br.com.dillmann.nginxignition.api.common.routing.RouteProvider
import org.koin.core.module.Module
import org.koin.dsl.bind

internal fun Module.frontendBeans() {
    single { FrontendRequestHandler() } bind RequestHandler::class
    single { FrontendRoutes(get()) } bind RouteProvider::class
}
