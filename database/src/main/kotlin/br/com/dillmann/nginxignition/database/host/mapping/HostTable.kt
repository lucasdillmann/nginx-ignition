package br.com.dillmann.nginxignition.database.host.mapping

import org.jetbrains.exposed.sql.Table

internal object HostTable: Table("host") {
    val id = uuid("id")
    val enabled = bool("enabled")
    val useGlobalBindings = bool("use_global_bindings")
    val defaultServer = bool("default_server")
    val domainNames = array<String>("domain_names").nullable()
    val websocketSupport = bool("websocket_support")
    val http2Support = bool("http2_support")
    val redirectHttpToHttps = bool("redirect_http_to_https")
    val accessListId = uuid("access_list_id").nullable()

    override val primaryKey = PrimaryKey(id)
}
