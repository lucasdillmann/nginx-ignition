package br.com.dillmann.nginxignition.database.settings.mapping

import org.jetbrains.exposed.sql.Table

@Suppress("MagicNumber")
internal object SettingsLogRotationTable: Table("settings_log_rotation") {
    val id = uuid("id")
    val enabled = bool("enabled")
    val maximumLines = integer("maximum_lines")
    val intervalUnit = varchar("interval_unit", 32)
    val intervalUnitCount = integer("interval_unit_count")

    override val primaryKey = PrimaryKey(id)
}
