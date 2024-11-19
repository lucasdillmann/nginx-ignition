package br.com.dillmann.nginxignition.api.certificate.model

import br.com.dillmann.nginxignition.api.common.serialization.UuidString
import kotlinx.serialization.Serializable

@Serializable
internal data class IssueCertificateResponse(
    val success: Boolean,
    val errorReason: String? = null,
    val certificateId: UuidString? = null,
)
