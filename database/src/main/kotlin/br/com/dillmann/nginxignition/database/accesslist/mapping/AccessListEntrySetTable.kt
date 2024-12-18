package br.com.dillmann.nginxignition.database.accesslist.mapping

import org.jetbrains.exposed.sql.Table

@Suppress("MagicNumber")
internal object AccessListEntrySetTable: Table("access_list_entry_set") {
    val id = uuid("id")
    val accessListId = uuid("access_list_id") references AccessListTable.id
    val priority = integer("priority")
    val outcome = varchar("outcome", 8)
    val sourceAddresses = array<String>("source_addresses")

    override val primaryKey = PrimaryKey(id)
}
