package br.com.dillmann.nginxignition.core.certificate

import java.time.OffsetDateTime
import java.util.UUID

data class Certificate(
    val id: UUID,
    val hosts: List<String>,
    val providerId: String,
    val issuedAt: OffsetDateTime,
    val validUntil: OffsetDateTime,
    val validFrom: OffsetDateTime,
    val renewAfter: OffsetDateTime?,
    val privateKey: String,
    val publicKey: String,
    val certificationChain: String?,
    val metadata: Map<String, Any>,
)
