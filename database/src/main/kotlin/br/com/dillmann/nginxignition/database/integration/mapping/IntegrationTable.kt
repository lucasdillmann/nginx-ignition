package br.com.dillmann.nginxignition.database.integration.mapping

import org.jetbrains.exposed.sql.Table

@Suppress("MagicNumber")
internal object IntegrationTable: Table("integration") {
    val id = varchar("id", 128)
    val enabled = bool("enabled")
    val parameters = text("parameters")

    override val primaryKey = PrimaryKey(id)
}
