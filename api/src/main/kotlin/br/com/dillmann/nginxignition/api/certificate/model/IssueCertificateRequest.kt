package br.com.dillmann.nginxignition.api.certificate.model

import kotlinx.serialization.Serializable
import kotlinx.serialization.json.JsonObject

@Serializable
internal data class IssueCertificateRequest(
    val providerId: String,
    val domainNames: List<String>,
    val parameters: JsonObject,
)
