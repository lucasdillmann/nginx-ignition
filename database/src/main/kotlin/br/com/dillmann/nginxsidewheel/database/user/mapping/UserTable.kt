package br.com.dillmann.nginxsidewheel.database.user.mapping

import org.jetbrains.exposed.sql.Table

internal object UserTable: Table("user") {
    val id = uuid("id")
    val enabled = bool("enabled")
    val name = varchar("name", 256)
    val username = varchar("username", 256)
    val passwordHash = varchar("password_hash", 2048)
    val passwordSalt = varchar("password_salt", 512)
    val role = varchar("role", 32)

    override val primaryKey = PrimaryKey(id)
}
