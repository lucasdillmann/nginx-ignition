package br.com.dillmann.nginxsidewheel.database.host.mapping

import org.jetbrains.exposed.sql.Table

internal object HostRouteTable: Table() {
    val id = uuid("id")
    val hostId = uuid("host_id") references HostTable.id
    val priority = integer("priority")
    val type = varchar("type", 64)
    val sourcePath = varchar("source_path", 512)
    val targetUri = varchar("target_uri", 512).nullable()
    val customSettings = text("custom_settings").nullable()
    val staticResponseCode = integer("static_response_code").nullable()
    val staticResponsePayload = text("static_response_payload").nullable()
    val staticResponseHeaders = text("static_response_headers").nullable()

    override val primaryKey = PrimaryKey(id)
}
