package br.com.dillmann.nginxignition.database.accesslist

import br.com.dillmann.nginxignition.core.accesslist.AccessList
import br.com.dillmann.nginxignition.database.accesslist.mapping.AccessListCredentialsTable
import br.com.dillmann.nginxignition.database.accesslist.mapping.AccessListEntrySetTable
import br.com.dillmann.nginxignition.database.accesslist.mapping.AccessListTable
import org.jetbrains.exposed.sql.ResultRow
import org.jetbrains.exposed.sql.statements.InsertStatement
import java.util.*

internal class AccessListConverter {
    fun apply(accessList: AccessList, scope: InsertStatement<out Any>) {
        with(AccessListTable) {
            scope[id] = accessList.id
            scope[name] = accessList.name
            scope[realm] = accessList.realm
            scope[satisfyAll] = accessList.satisfyAll
            scope[defaultOutcome] = accessList.defaultOutcome.name
            scope[forwardAuthenticationHeader] = accessList.forwardAuthenticationHeader
        }
    }

    fun apply(parentId: UUID, credentials: AccessList.Credentials, scope: InsertStatement<out Any>) {
        with(AccessListCredentialsTable) {
            scope[id] = UUID.randomUUID()
            scope[accessListId] = parentId
            scope[username] = credentials.username
            scope[password] = credentials.password
        }
    }

    fun apply(parentId: UUID, entrySet: AccessList.EntrySet, scope: InsertStatement<out Any>) {
        with(AccessListEntrySetTable) {
            scope[id] = UUID.randomUUID()
            scope[accessListId] = parentId
            scope[priority] = entrySet.priority
            scope[outcome] = entrySet.outcome.name
            scope[sourceAddresses] = entrySet.sourceAddresses
        }
    }

    fun toAccessList(accessList: ResultRow, credentials: List<ResultRow>, entries: List<ResultRow>) =
        AccessList(
            id = accessList[AccessListTable.id],
            name = accessList[AccessListTable.name],
            realm = accessList[AccessListTable.realm],
            satisfyAll = accessList[AccessListTable.satisfyAll],
            defaultOutcome = accessList[AccessListTable.defaultOutcome].let(AccessList.Outcome::valueOf),
            forwardAuthenticationHeader = accessList[AccessListTable.forwardAuthenticationHeader],
            entries = entries.map(::toEntrySet),
            credentials = credentials.map(::toCredentials),
        )

    private fun toEntrySet(entrySet: ResultRow) =
        AccessList.EntrySet(
            priority = entrySet[AccessListEntrySetTable.priority],
            outcome = entrySet[AccessListEntrySetTable.outcome].let(AccessList.Outcome::valueOf),
            sourceAddresses = entrySet[AccessListEntrySetTable.sourceAddresses],
        )

    private fun toCredentials(credentials: ResultRow) =
        AccessList.Credentials(
            username = credentials[AccessListCredentialsTable.username],
            password = credentials[AccessListCredentialsTable.password],
        )
}
