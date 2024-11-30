package br.com.dillmann.nginxignition.core.integration.command

fun interface ListIntegrationsCommand {
    data class Output(
        val id: String,
        val imageId: String?,
        val name: String,
        val description: String,
        val enabled: Boolean,
    )

    suspend fun getIntegrations(): List<Output>
}
