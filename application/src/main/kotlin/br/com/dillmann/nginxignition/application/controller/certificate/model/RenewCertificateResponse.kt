package br.com.dillmann.nginxignition.application.controller.certificate.model

import kotlinx.serialization.Serializable

@Serializable
data class RenewCertificateResponse(
    val success: Boolean,
    val errorReason: String? = null,
)
