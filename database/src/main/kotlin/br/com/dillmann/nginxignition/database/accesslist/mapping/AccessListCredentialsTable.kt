package br.com.dillmann.nginxignition.database.accesslist.mapping

import org.jetbrains.exposed.sql.Table

@Suppress("MagicNumber")
internal object AccessListCredentialsTable: Table("access_list_credentials") {
    val id = uuid("id")
    val accessListId = uuid("access_list_id") references AccessListTable.id
    val username = varchar("username", 256)
    val password = varchar("password", 256)

    override val primaryKey = PrimaryKey(id)
}
