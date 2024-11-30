package br.com.dillmann.nginxignition.core.integration

import br.com.dillmann.nginxignition.core.common.pagination.Page
import br.com.dillmann.nginxignition.core.integration.command.ConfigureIntegrationByIdCommand
import br.com.dillmann.nginxignition.core.integration.command.GetIntegrationByIdCommand
import br.com.dillmann.nginxignition.core.integration.command.ListIntegrationOptionsCommand
import br.com.dillmann.nginxignition.core.integration.command.ListIntegrationsCommand

internal class IntegrationService(
    private val repository: IntegrationRepository,
    private val adapters: List<IntegrationAdapter>,
) : ListIntegrationsCommand, GetIntegrationByIdCommand, ListIntegrationOptionsCommand, ConfigureIntegrationByIdCommand {
    private companion object {
        private val DEFAULT_SETTINGS = Integration(id = "", enabled = false, parameters = emptyMap())
    }

    override suspend fun getIntegrations(): List<ListIntegrationsCommand.Output> =
        adapters.map {
            val settings = repository.findById(it.id) ?: DEFAULT_SETTINGS

            ListIntegrationsCommand.Output(
                id = it.id,
                imageId = it.imageId,
                name = it.name,
                description = it.description,
                enabled = settings.enabled,
            )
        }

    override suspend fun getIntegrationById(id: String): GetIntegrationByIdCommand.Output? {
        val adapter = findAdapter(id) ?: return null
        val settings = repository.findById(id) ?: DEFAULT_SETTINGS

        return GetIntegrationByIdCommand.Output(
            id = id,
            imageId = adapter.imageId,
            name = adapter.name,
            description = adapter.description,
            enabled = settings.enabled,
            parameters = settings.parameters,
            configurationFields = adapter.configurationFields,
        )
    }

    override suspend fun getIntegrationOptions(
        integrationId: String,
        pageNumber: Int,
        pageSize: Int
    ): Page<ListIntegrationOptionsCommand.Output> {
        val adapter =
            findAdapter(integrationId) ?: return Page.empty()
        val settings =
            repository.findById(integrationId) ?: error("Integration with ID $integrationId it not configured")
        return adapter
            .getAvailableOptions(settings.parameters, pageNumber, pageSize)
            .map {
                ListIntegrationOptionsCommand.Output(
                    id = it.id,
                    name = it.name,
                )
            }
    }

    override suspend fun configureIntegration(id: String, enabled: Boolean, parameters: Map<String, Any?>) {
        // TODO: Validate the values before saving
        val configuration = Integration(id, enabled, parameters)
        repository.save(configuration)
    }

    suspend fun getIntegrationOptionUrl(
        integrationId: String,
        optionId: String,
    ): String {
        val adapter =
            findAdapter(integrationId) ?: error("No integration with ID $integrationId found")
        val settings =
            repository.findById(integrationId) ?: error("Integration with ID $integrationId it not configured")
        return adapter.getOptionProxyUrl(optionId, settings.parameters)
    }

    private fun findAdapter(id: String): IntegrationAdapter? =
        adapters.firstOrNull { it.id == id }
}
