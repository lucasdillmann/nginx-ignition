package br.com.dillmann.nginxignition.api.certificate.model

import br.com.dillmann.nginxignition.api.common.dynamicfield.DynamicFieldResponse
import kotlinx.serialization.Serializable

@Serializable
internal data class AvailableProviderResponse(
    val id: String,
    val name: String,
    val dynamicFields: List<DynamicFieldResponse>,
)
