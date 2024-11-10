package br.com.dillmann.nginxignition.database.certificate

import br.com.dillmann.nginxignition.core.certificate.Certificate
import br.com.dillmann.nginxignition.database.certificate.mapping.CertificateTable
import br.com.dillmann.nginxignition.database.common.json.toJsonObject
import br.com.dillmann.nginxignition.database.common.json.toPlainMap
import org.jetbrains.exposed.sql.ResultRow
import org.jetbrains.exposed.sql.statements.UpsertStatement

internal class CertificateConverter {
    fun apply(certificate: Certificate, scope: UpsertStatement<out Any>) {
        with(CertificateTable) {
            scope[id] = certificate.id
            scope[domainNames] = certificate.domainNames
            scope[providerId] = certificate.providerId
            scope[issuedAt] = certificate.issuedAt
            scope[validUntil] = certificate.validUntil
            scope[validFrom] = certificate.validFrom
            scope[renewAfter] = certificate.renewAfter
            scope[privateKey] = certificate.privateKey
            scope[publicKey] = certificate.publicKey
            scope[certificationChain] = certificate.certificationChain
            scope[parameters] = certificate.parameters.toJsonObject().toString()
            scope[metadata] = certificate.metadata
        }
    }

    fun toDomainModel(certificate: ResultRow) =
        with(CertificateTable) {
            Certificate(
                id = certificate[id],
                domainNames = certificate[domainNames],
                providerId = certificate[providerId],
                issuedAt = certificate[issuedAt],
                validUntil = certificate[validUntil],
                validFrom = certificate[validFrom],
                renewAfter = certificate[renewAfter],
                privateKey = certificate[privateKey],
                publicKey = certificate[publicKey],
                certificationChain = certificate[certificationChain],
                parameters = certificate[parameters].toJsonObject().toPlainMap(),
                metadata = certificate[metadata],
            )
        }
}
