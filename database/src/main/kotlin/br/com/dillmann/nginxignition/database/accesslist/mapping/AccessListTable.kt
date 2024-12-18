package br.com.dillmann.nginxignition.database.accesslist.mapping

import org.jetbrains.exposed.sql.Table

@Suppress("MagicNumber")
internal object AccessListTable: Table("access_list") {
    val id = uuid("id")
    val name = varchar("name", 256)
    val realm = varchar("realm", 256).nullable()
    val defaultOutcome = varchar("default_outcome", 8)
    val forwardAuthenticationHeader = bool("forward_authentication_header")

    override val primaryKey = PrimaryKey(id)
}
