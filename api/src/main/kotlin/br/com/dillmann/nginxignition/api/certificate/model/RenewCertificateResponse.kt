package br.com.dillmann.nginxignition.api.certificate.model

import kotlinx.serialization.Serializable

@Serializable
internal data class RenewCertificateResponse(
    val success: Boolean,
    val errorReason: String? = null,
)
