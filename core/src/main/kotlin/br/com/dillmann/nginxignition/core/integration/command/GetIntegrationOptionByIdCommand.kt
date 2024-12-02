package br.com.dillmann.nginxignition.core.integration.command

import br.com.dillmann.nginxignition.core.integration.model.IntegrationOption

fun interface GetIntegrationOptionByIdCommand {
    suspend fun getIntegrationOptionById(integrationId: String, optionId: String): IntegrationOption?
}
