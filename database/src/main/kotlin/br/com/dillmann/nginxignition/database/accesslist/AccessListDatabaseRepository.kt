package br.com.dillmann.nginxignition.database.accesslist

import br.com.dillmann.nginxignition.core.accesslist.AccessList
import br.com.dillmann.nginxignition.core.accesslist.AccessListRepository
import br.com.dillmann.nginxignition.core.common.pagination.Page
import br.com.dillmann.nginxignition.database.accesslist.mapping.AccessListCredentialsTable
import br.com.dillmann.nginxignition.database.accesslist.mapping.AccessListEntrySetTable
import br.com.dillmann.nginxignition.database.accesslist.mapping.AccessListTable
import br.com.dillmann.nginxignition.database.common.transaction.coTransaction
import br.com.dillmann.nginxignition.database.common.withSearchTerms
import org.jetbrains.exposed.sql.*
import org.jetbrains.exposed.sql.SqlExpressionBuilder.eq
import java.util.*

internal class AccessListDatabaseRepository(private val converter: AccessListConverter): AccessListRepository {
    data class Related(
        val credentials: List<ResultRow>,
        val entries: List<ResultRow>,
    )

    override suspend fun findById(id: UUID): AccessList? =
        findOneWhere { AccessListTable.id eq id }

    override suspend fun findByName(name: String): AccessList? =
        findOneWhere { AccessListTable.name.lowerCase() eq name.lowercase() }

    override suspend fun deleteById(id: UUID) {
        coTransaction {
            cleanupById(id)
        }
    }

    override suspend fun save(accessList: AccessList) {
        coTransaction {
            cleanupById(accessList.id)

            AccessListTable.insert {
                converter.apply(accessList, it)
            }

            accessList.credentials.forEach { credentials ->
                AccessListCredentialsTable.insert {
                    converter.apply(accessList.id, credentials, it)
                }
            }

            accessList.entries.forEach { entrySet ->
                AccessListEntrySetTable.insert {
                    converter.apply(accessList.id, entrySet, it)
                }
            }
        }
    }

    override suspend fun findPage(pageSize: Int, pageNumber: Int, searchTerms: String?): Page<AccessList> =
        coTransaction {
            val totalCount = AccessListTable
                .select(AccessListTable.id)
                .withSearchTerms(searchTerms, listOf(AccessListTable.name, AccessListTable.realm))
                .count()
            val accessLists = findAccessLists(pageSize, pageNumber, searchTerms)

            Page(
                contents = accessLists,
                pageNumber = pageNumber,
                pageSize = pageSize,
                totalItems = totalCount,
            )
        }

    override suspend fun findAll(): List<AccessList> =
        findAccessLists(null, null, null, null)

    private suspend fun findOneWhere(expression: SqlExpressionBuilder.() -> Op<Boolean>): AccessList? =
        coTransaction {
            val accessList = AccessListTable
                .select(AccessListTable.fields)
                .where(expression)
                .singleOrNull()
                ?: return@coTransaction null

            val (bindings, routes) = findRelated(accessList[AccessListTable.id])
            converter.toAccessList(accessList, bindings, routes)
        }

    private fun findAccessLists(
        pageSize: Int?,
        pageNumber: Int?,
        searchTerms: String?,
        predicate: (SqlExpressionBuilder.() -> Op<Boolean>)? = null,
    ): List<AccessList> =
        AccessListTable
            .select(AccessListTable.fields)
            .withSearchTerms(searchTerms, listOf(AccessListTable.name, AccessListTable.realm))
            .also {
                if (pageSize != null && pageNumber != null) {
                    it.limit(pageSize, pageSize.toLong() * pageNumber)
                }

                if (predicate != null) {
                    it.where(predicate)
                }
            }
            .orderBy(AccessListTable.name)
            .toList()
            .map { accessList ->
                val id = accessList[AccessListTable.id]
                val (credentials, entries) = findRelated(id)
                converter.toAccessList(accessList, credentials, entries)
            }

    private fun findRelated(id: UUID): Related {
        val bindings = AccessListCredentialsTable
            .select(AccessListCredentialsTable.fields)
            .where { AccessListCredentialsTable.accessListId eq id }
            .toList()

        val routes = AccessListEntrySetTable
            .select(AccessListEntrySetTable.fields)
            .where { AccessListEntrySetTable.accessListId eq id }
            .toList()

        return Related(bindings, routes)
    }

    private fun cleanupById(id: UUID) {
        AccessListEntrySetTable.deleteWhere { accessListId eq id }
        AccessListCredentialsTable.deleteWhere { accessListId eq id }
        AccessListTable.deleteWhere { AccessListTable.id eq id }
    }
}
