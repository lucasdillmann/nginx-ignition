package br.com.dillmann.nginxignition.api.integration.model

import br.com.dillmann.nginxignition.api.common.dynamicfield.DynamicFieldResponse
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.JsonObject

@Serializable
internal data class IntegrationConfigurationResponse(
    val id: String,
    val name: String,
    val description: String,
    val enabled: Boolean,
    val configurationFields: List<DynamicFieldResponse>,
    val parameters: JsonObject,
)
