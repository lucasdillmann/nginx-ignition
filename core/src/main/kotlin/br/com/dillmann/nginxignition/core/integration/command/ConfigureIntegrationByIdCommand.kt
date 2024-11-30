package br.com.dillmann.nginxignition.core.integration.command

fun interface ConfigureIntegrationByIdCommand {
    suspend fun configureIntegration(id: String, enabled: Boolean, parameters: Map<String, Any?>)
}
