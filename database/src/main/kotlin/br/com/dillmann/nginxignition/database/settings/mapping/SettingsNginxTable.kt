package br.com.dillmann.nginxignition.database.settings.mapping

import org.jetbrains.exposed.sql.Table

@Suppress("MagicNumber")
internal object SettingsNginxTable: Table("settings_nginx") {
    val id = uuid("id")
    val workerProcesses = integer("worker_processes")
    val serverTokensEnabled = bool("server_tokens_enabled")
    val readTimeout = integer("read_timeout")
    val connectTimeout = integer("connect_timeout")
    val sendTimeout = integer("send_timeout")
    val serverLogsEnabled = bool("server_logs_enabled")
    val accessLogsEnabled = bool("access_logs_enabled")
    val accessLogsFormat = varchar("access_logs_format", 512).nullable()
    val errorLogsEnabled = bool("error_logs_enabled")
    val errorLogsFormat = varchar("error_logs_format", 512).nullable()

    override val primaryKey = PrimaryKey(id)
}
