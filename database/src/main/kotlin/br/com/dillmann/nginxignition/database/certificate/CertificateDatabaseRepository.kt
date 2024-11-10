package br.com.dillmann.nginxignition.database.certificate

import br.com.dillmann.nginxignition.core.certificate.Certificate
import br.com.dillmann.nginxignition.core.certificate.CertificateRepository
import br.com.dillmann.nginxignition.core.common.pagination.Page
import br.com.dillmann.nginxignition.database.certificate.mapping.CertificateTable
import br.com.dillmann.nginxignition.database.common.transaction.coTransaction
import org.jetbrains.exposed.sql.SqlExpressionBuilder.eq
import org.jetbrains.exposed.sql.deleteWhere
import org.jetbrains.exposed.sql.upsert
import java.util.*

internal class CertificateDatabaseRepository(private val converter: CertificateConverter): CertificateRepository {
    override suspend fun findById(id: UUID): Certificate? =
        coTransaction {
            val certificate = CertificateTable
                .select(CertificateTable.fields)
                .where { CertificateTable.id eq id }
                .singleOrNull()
                ?: return@coTransaction null

            converter.toDomainModel(certificate)
        }

    override suspend fun deleteById(id: UUID) {
        coTransaction {
            CertificateTable.deleteWhere { CertificateTable.id eq id }
        }
    }

    override suspend fun save(certificate: Certificate) {
        coTransaction {
            CertificateTable.upsert { converter.apply(certificate, it) }
        }
    }

    override suspend fun findPage(pageSize: Int, pageNumber: Int): Page<Certificate> =
        coTransaction {
            val totalCount = CertificateTable.select(CertificateTable.id).count()
            val certificates = CertificateTable
                .select(CertificateTable.columns)
                .limit(pageSize, pageSize.toLong() * pageNumber)
                .map { converter.toDomainModel(it) }

            Page(
                contents = certificates,
                pageNumber = pageNumber,
                pageSize = pageSize,
                totalItems = totalCount,
            )
        }
}
