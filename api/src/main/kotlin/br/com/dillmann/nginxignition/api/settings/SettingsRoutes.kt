package br.com.dillmann.nginxignition.api.settings

import br.com.dillmann.nginxignition.api.common.routing.RouteNode
import br.com.dillmann.nginxignition.api.common.routing.RouteProvider
import br.com.dillmann.nginxignition.api.common.routing.routes
import br.com.dillmann.nginxignition.api.settings.handler.GetSettingsHandler
import br.com.dillmann.nginxignition.api.settings.handler.PutSettingsHandler
import br.com.dillmann.nginxignition.core.user.User

internal class SettingsRoutes(
    private val getHandler: GetSettingsHandler,
    private val putHandler: PutSettingsHandler,
): RouteProvider {
    override fun apiRoutes(): RouteNode =
        routes("/api/settings") {
            requireAuthentication {
                get(getHandler)

                requireRole(User.Role.ADMIN) {
                    put(putHandler)
                }
            }
        }
}
