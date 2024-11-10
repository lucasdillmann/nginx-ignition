package br.com.dillmann.nginxignition.application.controller.certificate.model

import br.com.dillmann.nginxignition.application.common.serialization.UuidString
import kotlinx.serialization.Serializable

@Serializable
data class IssueCertificateResponse(
    val success: Boolean,
    val errorReason: String? = null,
    val certificateId: UuidString? = null,
)
