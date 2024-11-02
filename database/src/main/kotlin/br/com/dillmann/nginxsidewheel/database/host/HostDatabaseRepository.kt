package br.com.dillmann.nginxsidewheel.database.host

import br.com.dillmann.nginxsidewheel.core.common.pagination.Page
import br.com.dillmann.nginxsidewheel.core.host.Host
import br.com.dillmann.nginxsidewheel.core.host.HostRepository
import br.com.dillmann.nginxsidewheel.database.host.mapping.HostBindingTable
import br.com.dillmann.nginxsidewheel.database.host.mapping.HostRouteTable
import br.com.dillmann.nginxsidewheel.database.host.mapping.HostTable
import org.jetbrains.exposed.sql.ResultRow
import org.jetbrains.exposed.sql.SqlExpressionBuilder.eq
import org.jetbrains.exposed.sql.deleteWhere
import org.jetbrains.exposed.sql.insert
import org.jetbrains.exposed.sql.transactions.transaction
import java.util.*

internal class HostDatabaseRepository(private val converter: HostConverter): HostRepository {
    data class Related(
        val bindings: List<ResultRow>,
        val routes: List<ResultRow>,
    )

    override suspend fun findById(id: UUID): Host? =
        transaction {
            val host = HostTable
                .select(HostTable.fields)
                .where { HostTable.id eq id }
                .singleOrNull()
                ?: return@transaction null

            val (bindings, routes) = findRelated(id)
            converter.toHost(host, bindings, routes)
        }

    override suspend fun deleteById(id: UUID) {
        transaction {
            cleanupById(id)
        }
    }

    override suspend fun save(host: Host) {
        transaction {
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
                HostRouteTable.insert {
                    converter.apply(host.id, binding, it)
                }
            }
        }
    }

    override suspend fun findPage(pageSize: Int, pageNumber: Int): Page<Host> =
        transaction {
            val totalCount = HostTable.select(HostTable.id).count()
            val hosts = findHosts(pageSize, pageNumber)

            Page(
                contents = hosts,
                pageNumber = pageNumber,
                pageSize = pageSize,
                totalItems = totalCount,
            )
        }

    override suspend fun findAll(): List<Host> =
        transaction {
            findHosts(null, null)
        }

    private fun findHosts(pageSize: Int?, pageNumber: Int?): List<Host> =
        HostTable
            .select(HostTable.fields)
            .also {
                if (pageSize != null && pageNumber != null) {
                    it.limit(pageSize, pageSize.toLong() * pageNumber)
                }
            }

            .orderBy(HostTable.id)
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
