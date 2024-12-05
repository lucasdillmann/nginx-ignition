package br.com.dillmann.nginxignition.database.certificate

import br.com.dillmann.nginxignition.core.certificate.Certificate
import br.com.dillmann.nginxignition.core.certificate.CertificateRepository
import br.com.dillmann.nginxignition.core.common.pagination.Page
import br.com.dillmann.nginxignition.database.certificate.mapping.CertificateTable
import br.com.dillmann.nginxignition.database.common.database.DatabaseState
import br.com.dillmann.nginxignition.database.common.database.DatabaseType
import br.com.dillmann.nginxignition.database.common.transaction.coTransaction
import br.com.dillmann.nginxignition.database.common.withSearchTerms
import org.jetbrains.exposed.sql.SqlExpressionBuilder.eq
import org.jetbrains.exposed.sql.and
import org.jetbrains.exposed.sql.deleteWhere
import org.jetbrains.exposed.sql.update
import org.jetbrains.exposed.sql.upsert
import java.time.OffsetDateTime
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

    override suspend fun existsById(id: UUID): Boolean =
        coTransaction {
            CertificateTable
                .select(CertificateTable.id)
                .where { CertificateTable.id eq id }
                .count() > 0
        }

    override suspend fun deleteById(id: UUID) {
        coTransaction {
            CertificateTable.deleteWhere { CertificateTable.id eq id }
        }
    }

    override suspend fun save(certificate: Certificate) {
        coTransaction {
            CertificateTable.upsert { converter.apply(certificate, it) }

            if (DatabaseState.type == DatabaseType.H2) {
                // H2 in the current version has a bug on upsert/merge operations where arrays values end up merged
                // as if they were only one (["a", "b"] becomes ["a, b"]). This update fixes that while H2 doesn't fix
                // it in the DBMS itself.
                CertificateTable.update(
                    where = { CertificateTable.id eq certificate.id },
                    body = {
                        it[domainNames] = certificate.domainNames
                        it[certificationChain] = certificate.certificationChain
                    }
                )
            }
        }
    }

    override suspend fun findPage(pageSize: Int, pageNumber: Int, searchTerms: String?): Page<Certificate> =
        coTransaction {
            val totalCount = CertificateTable
                .select(CertificateTable.id)
                .withSearchTerms(searchTerms, listOf(CertificateTable.domainNames))
                .count()
            val certificates = CertificateTable
                .select(CertificateTable.fields)
                .withSearchTerms(searchTerms, listOf(CertificateTable.domainNames))
                .limit(pageSize, pageSize.toLong() * pageNumber)
                .orderBy(CertificateTable.domainNames)
                .map { converter.toDomainModel(it) }

            Page(
                contents = certificates,
                pageNumber = pageNumber,
                pageSize = pageSize,
                totalItems = totalCount,
            )
        }

    override suspend fun findAllDueToRenew(): List<Certificate> =
        coTransaction {
            CertificateTable
                .select(CertificateTable.fields)
                .where {
                    val renewAfterNotNull = CertificateTable.renewAfter neq null
                    val renewAfterInThePast = CertificateTable.renewAfter lessEq OffsetDateTime.now()
                    renewAfterNotNull and renewAfterInThePast
                }
                .map { converter.toDomainModel(it) }
        }
}
