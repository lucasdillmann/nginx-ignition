package br.com.dillmann.nginxignition.api.settings.model

import br.com.dillmann.nginxignition.api.host.model.HostRequest
import br.com.dillmann.nginxignition.core.settings.Settings
import kotlinx.serialization.Serializable

@Serializable
internal data class SettingsDto(
    val nginx: NginxSettings,
    val logRotation: LogRotation,
    val certificateAutoRenew: CertificateAutoRenew,
    val globalBindings: List<HostRequest.Binding>,
) {
    @Serializable
    data class LogRotation(
        val enabled: Boolean,
        val maximumLines: Int,
        val intervalUnit: Settings.TimeUnit,
        val intervalUnitCount: Int,
    )

    @Serializable
    data class CertificateAutoRenew(
        val enabled: Boolean,
        val intervalUnit: Settings.TimeUnit,
        val intervalUnitCount: Int,
    )

    @Serializable
    data class NginxSettings(
        val logs: NginxLogs,
        val timeouts: NginxTimeouts,
        val workerProcesses: Int,
        val workerConnections: Int,
        val serverTokensEnabled: Boolean,
        val sendfileEnabled: Boolean,
        val gzipEnabled: Boolean,
        val defaultContentType: String,
        val maximumBodySizeMb: Int,
    )

    @Serializable
    data class NginxTimeouts(
        val read: Int,
        val connect: Int,
        val send: Int,
        val keepalive: Int,
    )

    @Serializable
    data class NginxLogs(
        val serverLogsEnabled: Boolean,
        val serverLogsLevel: Settings.LogLevel,
        val accessLogsEnabled: Boolean,
        val errorLogsEnabled: Boolean,
        val errorLogsLevel: Settings.LogLevel,
    )
}
