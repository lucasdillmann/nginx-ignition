package br.com.dillmann.nginxignition.database.settings

import br.com.dillmann.nginxignition.core.host.Host
import br.com.dillmann.nginxignition.core.settings.Settings
import br.com.dillmann.nginxignition.database.settings.mapping.SettingsCertificateAutoRenewTable
import br.com.dillmann.nginxignition.database.settings.mapping.SettingsGlobalBindingTable
import br.com.dillmann.nginxignition.database.settings.mapping.SettingsLogRotationTable
import br.com.dillmann.nginxignition.database.settings.mapping.SettingsNginxTable
import org.jetbrains.exposed.sql.ResultRow
import org.jetbrains.exposed.sql.statements.InsertStatement
import org.jetbrains.exposed.sql.statements.UpdateStatement

internal class SettingsConverter {
    fun apply(settings: Settings.NginxSettings, scope: UpdateStatement) {
        with(SettingsNginxTable) {
            scope[workerProcesses] = settings.workerProcesses
            scope[workerConnections] = settings.workerConnections
            scope[serverTokensEnabled] = settings.serverTokensEnabled
            scope[sendfileEnabled] = settings.sendfileEnabled
            scope[gzipEnabled] = settings.gzipEnabled
            scope[defaultContentType] = settings.defaultContentType
            scope[maximumBodySizeMb] = settings.maximumBodySizeMb
            scope[readTimeout] = settings.timeouts.read
            scope[connectTimeout] = settings.timeouts.connect
            scope[sendTimeout] = settings.timeouts.send
            scope[keepaliveTimeout] = settings.timeouts.keepalive
            scope[serverLogsEnabled] = settings.logs.serverLogsEnabled
            scope[serverLogsLevel] = settings.logs.serverLogsLevel.name
            scope[accessLogsEnabled] = settings.logs.accessLogsEnabled
            scope[errorLogsEnabled] = settings.logs.errorLogsEnabled
            scope[errorLogsLevel] = settings.logs.errorLogsLevel.name
        }
    }

    fun apply(settings: Settings.LogRotation, scope: UpdateStatement) {
        with(SettingsLogRotationTable) {
            scope[enabled] = settings.enabled
            scope[maximumLines] = settings.maximumLines
            scope[intervalUnit] = settings.intervalUnit.name
            scope[intervalUnitCount] = settings.intervalUnitCount
        }
    }

    fun apply(settings: Settings.CertificateAutoRenew, scope: UpdateStatement) {
        with(SettingsCertificateAutoRenewTable) {
            scope[enabled] = settings.enabled
            scope[intervalUnit] = settings.intervalUnit.name
            scope[intervalUnitCount] = settings.intervalUnitCount
        }
    }

    fun apply(settings: Host.Binding, scope: InsertStatement<out Any>) {
        with(SettingsGlobalBindingTable) {
            scope[id] = settings.id
            scope[type] = settings.type.name
            scope[ip] = settings.ip
            scope[port] = settings.port
            scope[certificateId] = settings.certificateId
        }
    }

    fun toNginxSettings(settings: ResultRow) =
        with(SettingsNginxTable) {
            Settings.NginxSettings(
                workerProcesses = settings[workerProcesses],
                workerConnections = settings[workerConnections],
                serverTokensEnabled = settings[serverTokensEnabled],
                maximumBodySizeMb = settings[maximumBodySizeMb],
                defaultContentType = settings[defaultContentType],
                sendfileEnabled = settings[sendfileEnabled],
                gzipEnabled = settings[gzipEnabled],
                timeouts = Settings.NginxTimeouts(
                    read = settings[readTimeout],
                    connect = settings[connectTimeout],
                    send = settings[sendTimeout],
                    keepalive = settings[keepaliveTimeout],
                ),
                logs = Settings.NginxLogs(
                    serverLogsEnabled = settings[serverLogsEnabled],
                    serverLogsLevel = settings[serverLogsLevel].let(Settings.LogLevel::valueOf),
                    accessLogsEnabled = settings[accessLogsEnabled],
                    errorLogsEnabled = settings[errorLogsEnabled],
                    errorLogsLevel = settings[errorLogsLevel].let(Settings.LogLevel::valueOf),
                ),
            )
        }

    fun toGlobalHostBinding(settings: ResultRow) =
        with(SettingsGlobalBindingTable) {
            Host.Binding(
                id = settings[id],
                type = settings[type].let(Host.BindingType::valueOf),
                ip = settings[ip],
                port = settings[port],
                certificateId = settings[certificateId],
            )
        }

    fun toLogRotationSettings(settings: ResultRow) =
        with(SettingsLogRotationTable) {
            Settings.LogRotation(
                enabled = settings[enabled],
                maximumLines = settings[maximumLines],
                intervalUnit = settings[intervalUnit].let(Settings.TimeUnit::valueOf),
                intervalUnitCount = settings[intervalUnitCount],
            )
        }

    fun toCertificateAutoRenewSettings(settings: ResultRow) =
        with(SettingsCertificateAutoRenewTable) {
            Settings.CertificateAutoRenew(
                enabled = settings[enabled],
                intervalUnit = settings[intervalUnit].let(Settings.TimeUnit::valueOf),
                intervalUnitCount = settings[intervalUnitCount],
            )
        }
}
