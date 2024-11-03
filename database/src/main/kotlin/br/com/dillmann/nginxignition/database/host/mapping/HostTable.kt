package br.com.dillmann.nginxignition.database.host.mapping

import org.jetbrains.exposed.sql.Table

internal object HostTable: Table("host") {
    val id = uuid("id")
    val enabled = bool("enabled")
    val default = bool("default")
    val domainNames = varchar("domainNames", 512)
    val websocketSupport = bool("websocket_support")
    val http2Support = bool("http2_support")
    val redirectHttpToHttps = bool("redirect_http_to_https")

    override val primaryKey = PrimaryKey(id)
}
