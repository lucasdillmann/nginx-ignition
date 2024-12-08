package br.com.dillmann.nginxignition.core.settings

import br.com.dillmann.nginxignition.core.host.Host

data class Settings(
    val nginx: NginxSettings,
    val logRotation: LogRotation,
    val certificateAutoRenew: CertificateAutoRenew,
    val globalBindings: List<Host.Binding>,
) {
    data class LogRotation(
        val enabled: Boolean,
        val maximumLines: Int,
        val intervalUnit: TimeUnit,
        val intervalUnitCount: Int,
    )

    data class CertificateAutoRenew(
        val enabled: Boolean,
        val intervalUnit: TimeUnit,
        val intervalUnitCount: Int,
    )

    data class NginxSettings(
        val logs: NginxLogs,
        val timeouts: NginxTimeouts,
        val workerProcesses: Int,
        val workerConnections: Int,
        val defaultContentType: String,
        val serverTokensEnabled: Boolean,
        val maximumBodySizeMb: Int,
        val sendfileEnabled: Boolean,
        val gzipEnabled: Boolean,
    )

    data class NginxTimeouts(
        val read: Int,
        val connect: Int,
        val send: Int,
        val keepalive: Int,
    )

    data class NginxLogs(
        val serverLogsEnabled: Boolean,
        val serverLogsLevel: LogLevel,
        val accessLogsEnabled: Boolean,
        val errorLogsEnabled: Boolean,
        val errorLogsLevel: LogLevel,
    )

    enum class LogLevel {
        WARN,
        ERROR,
        CRIT,
        ALERT,
        EMERG,
    }

    enum class TimeUnit {
        MINUTES,
        HOURS,
        DAYS,
    }
}
