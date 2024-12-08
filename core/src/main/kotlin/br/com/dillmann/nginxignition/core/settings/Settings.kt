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
        val serverTokensEnabled: Boolean,
    )

    data class NginxTimeouts(
        val read: Int,
        val connect: Int,
        val send: Int,
    )

    data class NginxLogs(
        val serverLogsEnabled: Boolean,
        val accessLogsEnabled: Boolean,
        val accessLogsFormat: String?,
        val errorLogsEnabled: Boolean,
        val errorLogsFormat: String?,
    )

    enum class TimeUnit {
        MINUTES,
        HOURS,
        DAYS,
    }
}
