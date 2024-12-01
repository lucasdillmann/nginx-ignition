package br.com.dillmann.nginxignition.api.integration.model

import kotlinx.serialization.Serializable

@Serializable
internal data class IntegrationResponse(
    val id: String,
    val name: String,
    val description: String,
    val enabled: Boolean,
)
