import { HostBinding } from "../../host/model/HostRequest"

export enum TimeUnit {
    MINUTES = "MINUTES",
    HOURS = "HOURS",
    DAYS = "DAYS",
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
}

export interface NginxLogsSettingsDto {
    serverLogsEnabled: boolean
    accessLogsEnabled: boolean
    accessLogsFormat?: string
    errorLogsEnabled: boolean
    errorLogsFormat?: string
}

export interface NginxSettingsDto {
    logs: NginxLogsSettingsDto
    timeouts: NginxTimeoutsSettingsDto
    workerProcesses: number
    serverTokensEnabled: boolean
}

export default interface SettingsDto {
    nginx: NginxSettingsDto
    logRotation: LogRotationSettingsDto
    certificateAutoRenew: CertificateAutoRenewSettingsDto
    globalBindings: HostBinding[]
}
