package br.com.dillmann.nginxignition.api.integration.model

import kotlinx.serialization.Serializable

@Serializable
internal data class IntegrationOptionResponse(
    val id: String,
    val name: String,
)
