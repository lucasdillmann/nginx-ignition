package br.com.dillmann.nginxignition.database.host

import br.com.dillmann.nginxignition.core.common.pagination.Page
import br.com.dillmann.nginxignition.core.host.Host
import br.com.dillmann.nginxignition.core.host.HostRepository
import br.com.dillmann.nginxignition.database.common.transaction.coTransaction
import br.com.dillmann.nginxignition.database.common.withSearchTerms
import br.com.dillmann.nginxignition.database.host.mapping.HostBindingTable
import br.com.dillmann.nginxignition.database.host.mapping.HostRouteTable
import br.com.dillmann.nginxignition.database.host.mapping.HostTable
import org.jetbrains.exposed.sql.*
import org.jetbrains.exposed.sql.SqlExpressionBuilder.eq
import java.util.*

internal class HostDatabaseRepository(private val converter: HostConverter): HostRepository {
    data class Related(
        val bindings: List<ResultRow>,
        val routes: List<ResultRow>,
    )

    override suspend fun findById(id: UUID): Host? =
        findOneWhere { HostTable.id eq id }

    override suspend fun deleteById(id: UUID) {
        coTransaction {
            cleanupById(id)
        }
    }

    override suspend fun findDefault(): Host? =
        findOneWhere { HostTable.defaultServer eq true }

    override suspend fun save(host: Host) {
        coTransaction {
            cleanupById(host.id)

            HostTable.insert {
                converter.apply(host, it)
            }

            host.routes.forEach { route ->
                HostRouteTable.insert {
                    converter.apply(host.id, route, it)
                }
            }

            host.bindings.forEach { binding ->
                HostBindingTable.insert {
                    converter.apply(host.id, binding, it)
                }
            }
        }
    }

    override suspend fun findPage(pageSize: Int, pageNumber: Int, searchTerms: String?): Page<Host> =
        coTransaction {
            val totalCount = HostTable
                .select(HostTable.id)
                .withSearchTerms(searchTerms, listOf(HostTable.domainNames))
                .count()
            val hosts = findHosts(pageSize, pageNumber, searchTerms)

            Page(
                contents = hosts,
                pageNumber = pageNumber,
                pageSize = pageSize,
                totalItems = totalCount,
            )
        }

    override suspend fun findAllEnabled(): List<Host> =
        coTransaction {
            findHosts(null, null, null) { HostTable.enabled eq true }
        }

    override suspend fun existsById(id: UUID): Boolean =
        coTransaction {
            HostTable.select(HostTable.id).where { HostTable.id eq id }.count() > 0
        }

    override suspend fun existsByAccessListId(accessListId: UUID): Boolean =
        coTransaction {
            val hostExists = HostTable
                .select(HostTable.id)
                .where { HostTable.accessListId eq accessListId }
                .count() > 0
            if (hostExists) return@coTransaction true

            HostRouteTable
                .select(HostRouteTable.id)
                .where { HostRouteTable.accessListId eq accessListId }
                .count() > 0
        }

    override suspend fun existsByCertificateId(certificateId: UUID): Boolean =
        coTransaction {
            HostBindingTable
                .select(HostBindingTable.id)
                .where { HostBindingTable.certificateId eq certificateId }
                .count() > 0
        }

    private suspend fun findOneWhere(expression: SqlExpressionBuilder.() -> Op<Boolean>): Host? =
        coTransaction {
            val host = HostTable
                .select(HostTable.fields)
                .where(expression)
                .singleOrNull()
                ?: return@coTransaction null

            val (bindings, routes) = findRelated(host[HostTable.id])
            converter.toHost(host, bindings, routes)
        }

    private fun findHosts(
        pageSize: Int?,
        pageNumber: Int?,
        searchTerms: String?,
        predicate: (SqlExpressionBuilder.() -> Op<Boolean>)? = null,
    ): List<Host> =
        HostTable
            .select(HostTable.fields)
            .withSearchTerms(searchTerms, listOf(HostTable.domainNames))
            .also {
                if (pageSize != null && pageNumber != null) {
                    it.limit(pageSize).offset(pageSize.toLong() * pageNumber)
                }

                if (predicate != null) {
                    it.where(predicate)
                }
            }
            .orderBy(HostTable.domainNames)
            .toList()
            .map { host ->
                val id = host[HostTable.id]
                val (bindings, routes) = findRelated(id)
                converter.toHost(host, bindings, routes)
            }

    private fun findRelated(id: UUID): Related {
        val bindings = HostBindingTable
            .select(HostBindingTable.fields)
            .where { HostBindingTable.hostId eq id }
            .toList()

        val routes = HostRouteTable
            .select(HostRouteTable.fields)
            .where { HostRouteTable.hostId eq id }
            .toList()

        return Related(bindings, routes)
    }

    private fun cleanupById(id: UUID) {
        HostRouteTable.deleteWhere { hostId eq id }
        HostBindingTable.deleteWhere { hostId eq id }
        HostTable.deleteWhere { HostTable.id eq id }
    }
}
