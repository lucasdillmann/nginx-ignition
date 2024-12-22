package br.com.dillmann.nginxignition.api.integration

import br.com.dillmann.nginxignition.api.common.routing.*
import br.com.dillmann.nginxignition.api.integration.handler.GetIntegrationConfigurationHandler
import br.com.dillmann.nginxignition.api.integration.handler.GetIntegrationOptionByIdHandler
import br.com.dillmann.nginxignition.api.integration.handler.GetIntegrationOptionsHandler
import br.com.dillmann.nginxignition.api.integration.handler.ListIntegrationsHandler
import br.com.dillmann.nginxignition.api.integration.handler.PutIntegrationConfigurationHandler
import br.com.dillmann.nginxignition.core.user.User

internal class IntegrationRoutes(
    private val listHandler: ListIntegrationsHandler,
    private val getConfigurationHandler: GetIntegrationConfigurationHandler,
    private val putConfigurationHandler: PutIntegrationConfigurationHandler,
    private val getOptionsHandler: GetIntegrationOptionsHandler,
    private val getOptionByIdHandler: GetIntegrationOptionByIdHandler,
): RouteProvider {
    override fun apiRoutes(): RouteNode =
        basePath("/api/integrations") {
            requireAuthentication {
                get(listHandler)

                path("/{id}") {
                    path("/options") {
                        get(getOptionsHandler)
                        get("/{optionId}", getOptionByIdHandler)
                    }

                    path("/configuration") {
                        requireRole(User.Role.ADMIN) {
                            get(getConfigurationHandler)
                            put(putConfigurationHandler)
                        }
                    }
                }
            }
        }
}
