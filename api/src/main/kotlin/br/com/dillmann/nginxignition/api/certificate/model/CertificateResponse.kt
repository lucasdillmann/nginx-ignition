package br.com.dillmann.nginxignition.api.certificate.model

import br.com.dillmann.nginxignition.api.common.serialization.OffsetDateTimeString
import br.com.dillmann.nginxignition.api.common.serialization.UuidString
import kotlinx.serialization.Serializable

@Serializable
internal data class CertificateResponse(
    val id: UuidString,
    val domainNames: List<String>,
    val providerId: String,
    val issuedAt: OffsetDateTimeString,
    val validUntil: OffsetDateTimeString,
    val validFrom: OffsetDateTimeString,
    val renewAfter: OffsetDateTimeString?,
)
