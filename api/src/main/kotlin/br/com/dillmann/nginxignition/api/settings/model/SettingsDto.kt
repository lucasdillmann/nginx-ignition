package br.com.dillmann.nginxignition.api.settings.model

import br.com.dillmann.nginxignition.api.host.model.HostResponse
import br.com.dillmann.nginxignition.core.settings.Settings
import kotlinx.serialization.Serializable

@Serializable
internal data class SettingsDto(
    val nginx: NginxSettings,
    val logRotation: LogRotation,
    val certificateAutoRenew: CertificateAutoRenew,
    val globalBindings: List<HostResponse.Binding>,
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
        val serverTokensEnabled: Boolean,
    )

    @Serializable
    data class NginxTimeouts(
        val read: Int,
        val connect: Int,
        val send: Int,
    )

    @Serializable
    data class NginxLogs(
        val serverLogsEnabled: Boolean,
        val accessLogsEnabled: Boolean,
        val accessLogsFormat: String? = null,
        val errorLogsEnabled: Boolean,
        val errorLogsFormat: String? = null,
    )
}
