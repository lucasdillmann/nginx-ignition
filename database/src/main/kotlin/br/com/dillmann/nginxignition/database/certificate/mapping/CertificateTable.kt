package br.com.dillmann.nginxignition.database.certificate.mapping

import org.jetbrains.exposed.sql.Table
import org.jetbrains.exposed.sql.javatime.timestampWithTimeZone

@Suppress("MagicNumber")
internal object CertificateTable: Table("certificate") {
    val id = uuid("id")
    val domainNames = array<String>("domain_names")
    val providerId = varchar("provider_id", 64)
    val issuedAt = timestampWithTimeZone("issued_at")
    val validUntil = timestampWithTimeZone("valid_until")
    val validFrom = timestampWithTimeZone("valid_from")
    val renewAfter = timestampWithTimeZone("renew_after").nullable()
    val privateKey = varchar("private_key", 2048)
    val publicKey = varchar("public_key", 2048)
    val certificationChain = array<String>("certification_chain")
    val parameters = text("parameters")
    val metadata = text("metadata").nullable()

    override val primaryKey = PrimaryKey(id)
}
