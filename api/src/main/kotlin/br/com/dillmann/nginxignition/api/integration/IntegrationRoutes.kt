package br.com.dillmann.nginxignition.api.integration

import br.com.dillmann.nginxignition.api.common.routing.*
import br.com.dillmann.nginxignition.api.integration.handler.GetIntegrationConfigurationHandler
import br.com.dillmann.nginxignition.api.integration.handler.GetIntegrationOptionsHandler
import br.com.dillmann.nginxignition.api.integration.handler.ListIntegrationsHandler
import br.com.dillmann.nginxignition.api.integration.handler.PutIntegrationConfigurationHandler

internal class IntegrationRoutes(
    private val listHandler: ListIntegrationsHandler,
    private val getConfigurationHandler: GetIntegrationConfigurationHandler,
    private val putConfigurationHandler: PutIntegrationConfigurationHandler,
    private val getOptionsHandler: GetIntegrationOptionsHandler,
): RouteProvider {
    override fun apiRoutes(): RouteNode =
        routes("/api/integrations") {
            requireAuthentication {
                get(listHandler)
                get("/{id}/configuration", getConfigurationHandler)
                put("/{id}/configuration", putConfigurationHandler)
                get("/{id}/options", getOptionsHandler)
            }
        }
}
