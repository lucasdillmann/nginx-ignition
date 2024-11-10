package br.com.dillmann.nginxignition.application.controller.certificate.model

import kotlinx.serialization.Serializable
import kotlinx.serialization.json.JsonObject

@Serializable
data class IssueCertificateRequest(
    val providerId: String,
    val domainNames: List<String>,
    val parameters: JsonObject,
)
