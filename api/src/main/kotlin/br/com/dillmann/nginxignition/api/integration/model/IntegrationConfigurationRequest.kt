package br.com.dillmann.nginxignition.api.integration.model

import kotlinx.serialization.Serializable
import kotlinx.serialization.json.JsonObject

@Serializable
internal data class IntegrationConfigurationRequest(
    val enabled: Boolean,
    val parameters: JsonObject,
)
