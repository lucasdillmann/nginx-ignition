package br.com.dillmann.nginxignition.application.frontend

import br.com.dillmann.nginxignition.api.common.routing.RouteNode
import br.com.dillmann.nginxignition.api.common.routing.RouteProvider
import br.com.dillmann.nginxignition.api.common.routing.basePath

internal class FrontendRoutes(private val handler: FrontendRequestHandler): RouteProvider {
    override fun apiRoutes(): RouteNode =
        basePath("/**") {
            get(handler)
        }
}
