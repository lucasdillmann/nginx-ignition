package br.com.dillmann.nginxignition.database.settings.mapping

import org.jetbrains.exposed.sql.Table

@Suppress("MagicNumber")
internal object SettingsCertificateAutoRenewTable: Table("settings_certificate_auto_renew") {
    val id = uuid("id")
    val enabled = bool("enabled")
    val intervalUnit = varchar("interval_unit", 32)
    val intervalUnitCount = integer("interval_unit_count")

    override val primaryKey = PrimaryKey(id)
}
