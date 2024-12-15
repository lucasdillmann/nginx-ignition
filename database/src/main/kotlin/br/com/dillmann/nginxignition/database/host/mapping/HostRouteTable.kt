package br.com.dillmann.nginxignition.database.host.mapping

import org.jetbrains.exposed.sql.Table

@Suppress("MagicNumber")
internal object HostRouteTable: Table("host_route") {
    val id = uuid("id")
    val hostId = uuid("host_id") references HostTable.id
    val priority = integer("priority")
    val type = varchar("type", 64)
    val sourcePath = varchar("source_path", 512)
    val targetUri = varchar("target_uri", 512).nullable()
    val customSettings = text("custom_settings").nullable()
    val redirectCode = integer("redirect_code").nullable()
    val staticResponseCode = integer("static_response_code").nullable()
    val staticResponsePayload = text("static_response_payload").nullable()
    val staticResponseHeaders = text("static_response_headers").nullable()
    val integrationId = varchar("integration_id", 128).nullable()
    val integrationOptionId = varchar("integration_option_id", 255).nullable()
    val includeForwardHeaders = bool("include_forward_headers")
    val proxySslServerName = bool("proxy_ssl_server_name")
    val keepOriginalDomainName = bool("keep_original_domain_name")
    val forwardQueryParams = bool("forward_query_params")

    override val primaryKey = PrimaryKey(id)
}
