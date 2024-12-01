package br.com.dillmann.nginxignition.core.integration

import br.com.dillmann.nginxignition.core.common.dynamicfield.DynamicFields
import br.com.dillmann.nginxignition.core.common.pagination.Page
import br.com.dillmann.nginxignition.core.integration.command.ConfigureIntegrationByIdCommand
import br.com.dillmann.nginxignition.core.integration.command.GetIntegrationByIdCommand
import br.com.dillmann.nginxignition.core.integration.command.ListIntegrationOptionsCommand
import br.com.dillmann.nginxignition.core.integration.command.ListIntegrationsCommand
import br.com.dillmann.nginxignition.core.integration.exception.IntegrationDisabledException
import br.com.dillmann.nginxignition.core.integration.exception.IntegrationNotConfiguredException
import br.com.dillmann.nginxignition.core.integration.exception.IntegrationNotFoundException

internal class IntegrationService(
    private val repository: IntegrationRepository,
    private val adapters: List<IntegrationAdapter>,
    private val validator: IntegrationValidator,
) : ListIntegrationsCommand, GetIntegrationByIdCommand, ListIntegrationOptionsCommand, ConfigureIntegrationByIdCommand {
    private companion object {
        private val DEFAULT_SETTINGS = Integration(id = "", enabled = false, parameters = emptyMap())
    }

    override suspend fun getIntegrations(): List<ListIntegrationsCommand.Output> =
        adapters.map {
            val settings = repository.findById(it.id) ?: DEFAULT_SETTINGS

            ListIntegrationsCommand.Output(
                id = it.id,
                name = it.name,
                description = it.description,
                enabled = settings.enabled,
            )
        }

    override suspend fun getIntegrationById(id: String): GetIntegrationByIdCommand.Output {
        val adapter = findAdapter(id)
        val settings = repository.findById(id) ?: DEFAULT_SETTINGS

        return GetIntegrationByIdCommand.Output(
            id = id,
            name = adapter.name,
            description = adapter.description,
            enabled = settings.enabled,
            parameters = DynamicFields.removeSensitiveParameters(adapter.configurationFields, settings.parameters),
            configurationFields = adapter.configurationFields,
        )
    }

    override suspend fun getIntegrationOptions(
        integrationId: String,
        pageNumber: Int,
        pageSize: Int
    ): Page<ListIntegrationOptionsCommand.Output> {
        val adapter = findAdapter(integrationId)
        val settings = findSettings(integrationId)
        if (!settings.enabled) throw IntegrationDisabledException()

        return adapter
            .getAvailableOptions(settings.parameters, pageNumber, pageSize)
            .map { ListIntegrationOptionsCommand.Output(it.id, it.name) }
    }

    override suspend fun configureIntegration(id: String, enabled: Boolean, parameters: Map<String, Any?>) {
        val adapter = findAdapter(id)
        if (enabled)
            validator.validate(adapter.configurationFields, parameters)

        val configuration = Integration(id, enabled, parameters)
        repository.save(configuration)
    }

    suspend fun getIntegrationOptionUrl(
        integrationId: String,
        optionId: String,
    ): String {
        val adapter = findAdapter(integrationId)
        val settings = findSettings(integrationId)
        if (!settings.enabled) throw IntegrationDisabledException()

        return adapter.getOptionProxyUrl(optionId, settings.parameters)
    }

    private fun findAdapter(id: String): IntegrationAdapter =
        adapters.firstOrNull { it.id == id } ?: throw IntegrationNotFoundException()

    private suspend fun findSettings(id: String): Integration =
        repository.findById(id) ?: throw IntegrationNotConfiguredException()
}
