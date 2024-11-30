package br.com.dillmann.nginxignition.core.integration.command

import br.com.dillmann.nginxignition.core.common.dynamicfield.DynamicField

fun interface GetIntegrationByIdCommand {
    data class Output(
        val id: String,
        val imageId: String?,
        val name: String,
        val description: String,
        val enabled: Boolean,
        val configurationFields: List<DynamicField>,
        val parameters: Map<String, Any?>,
    )

    suspend fun getIntegrationById(id: String): Output?
}
