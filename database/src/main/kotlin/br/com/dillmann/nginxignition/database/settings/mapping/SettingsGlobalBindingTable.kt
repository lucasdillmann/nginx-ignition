package br.com.dillmann.nginxignition.database.settings.mapping

import org.jetbrains.exposed.sql.Table

@Suppress("MagicNumber")
internal object SettingsGlobalBindingTable: Table("settings_global_binding") {
    val id = uuid("id")
    val type = varchar("type", 64)
    val ip = varchar("ip", 256)
    val port = integer("port")
    val certificateId = uuid("certificate_id").nullable()

    override val primaryKey = PrimaryKey(id)
}
