package br.com.dillmann.nginxignition.application.controller.certificate.model

import br.com.dillmann.nginxignition.application.common.serialization.OffsetDateTimeString
import br.com.dillmann.nginxignition.application.common.serialization.UuidString
import kotlinx.serialization.Serializable

@Serializable
data class CertificateResponse(
    val id: UuidString,
    val domainNames: List<String>,
    val providerId: String,
    val issuedAt: OffsetDateTimeString,
    val validUntil: OffsetDateTimeString,
    val validFrom: OffsetDateTimeString,
    val renewAfter: OffsetDateTimeString?,
)
