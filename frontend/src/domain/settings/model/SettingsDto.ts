import { HostBinding } from "../../host/model/HostRequest"

export enum TimeUnit {
    MINUTES = "MINUTES",
    HOURS = "HOURS",
    DAYS = "DAYS",
}

export enum LogLevel {
    WARN = "WARN",
    ERROR = "ERROR",
    CRIT = "CRIT",
    ALERT = "ALERT",
    EMERG = "EMERG",
}

export interface CertificateAutoRenewSettingsDto {
    enabled: boolean
    intervalUnit: TimeUnit
    intervalUnitCount: number
}

export interface LogRotationSettingsDto {
    enabled: boolean
    maximumLines: number
    intervalUnit: TimeUnit
    intervalUnitCount: number
}

export interface NginxTimeoutsSettingsDto {
    read: number
    connect: number
    send: number
    keepalive: number
}

export interface NginxLogsSettingsDto {
    serverLogsEnabled: boolean
    serverLogsLevel: LogLevel
    accessLogsEnabled: boolean
    errorLogsEnabled: boolean
    errorLogsLevel: LogLevel
}

export interface NginxSettingsDto {
    logs: NginxLogsSettingsDto
    timeouts: NginxTimeoutsSettingsDto
    workerProcesses: number
    workerConnections: number
    serverTokensEnabled: boolean
    sendfileEnabled: boolean
    gzipEnabled: boolean
    maximumBodySizeMb: number
    defaultContentType: string
}

export default interface SettingsDto {
    nginx: NginxSettingsDto
    logRotation: LogRotationSettingsDto
    certificateAutoRenew: CertificateAutoRenewSettingsDto
    globalBindings: HostBinding[]
}
