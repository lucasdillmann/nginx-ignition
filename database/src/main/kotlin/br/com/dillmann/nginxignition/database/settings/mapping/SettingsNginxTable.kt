package br.com.dillmann.nginxignition.database.settings.mapping

import org.jetbrains.exposed.sql.Table

@Suppress("MagicNumber")
internal object SettingsNginxTable: Table("settings_nginx") {
    val id = uuid("id")
    val workerProcesses = integer("worker_processes")
    val workerConnections = integer("worker_connections")
    val serverTokensEnabled = bool("server_tokens_enabled")
    val sendfileEnabled = bool("sendfile_enabled")
    val gzipEnabled = bool("gzip_enabled")
    val defaultContentType = varchar("default_content_type", 128)
    val maximumBodySizeMb = integer("maximum_body_size_mb")
    val readTimeout = integer("read_timeout")
    val connectTimeout = integer("connect_timeout")
    val sendTimeout = integer("send_timeout")
    val keepaliveTimeout = integer("keepalive_timeout")
    val serverLogsEnabled = bool("server_logs_enabled")
    val serverLogsLevel = varchar("server_logs_level", 8)
    val accessLogsEnabled = bool("access_logs_enabled")
    val errorLogsEnabled = bool("error_logs_enabled")
    val errorLogsLevel = varchar("error_logs_level", 8)

    override val primaryKey = PrimaryKey(id)
}
