package br.com.dillmann.nginxignition.core.certificate

import java.time.OffsetDateTime
import java.util.UUID

data class Certificate(
    val id: UUID,
    val domainNames: List<String>,
    val providerId: String,
    val issuedAt: OffsetDateTime,
    val validUntil: OffsetDateTime,
    val validFrom: OffsetDateTime,
    val renewAfter: OffsetDateTime?,
    val privateKey: String,
    val publicKey: String,
    val certificationChain: List<String>,
    val parameters: Map<String, Any?>,
    val metadata: String?,
)
