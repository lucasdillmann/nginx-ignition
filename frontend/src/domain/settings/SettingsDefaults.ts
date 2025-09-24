import SettingsFormValues from "./model/SettingsFormValues"
import { LogLevel, RuntimeUser, TimeUnit } from "./model/SettingsDto"

export function settingsDefaults(): SettingsFormValues {
    return {
        nginx: {
            logs: {
                serverLogsEnabled: true,
                serverLogsLevel: LogLevel.ERROR,
                accessLogsEnabled: true,
                errorLogsEnabled: true,
                errorLogsLevel: LogLevel.ERROR,
            },
            timeouts: {
                connect: 5,
                keepalive: 30,
                read: 300,
                send: 300,
            },
            defaultContentType: "application/octet-stream",
            gzipEnabled: true,
            maximumBodySizeMb: 1024,
            sendfileEnabled: true,
            serverTokensEnabled: false,
            workerConnections: 1024,
            workerProcesses: 2,
            runtimeUser: RuntimeUser.ROOT,
        },
        logRotation: {
            enabled: true,
            intervalUnit: TimeUnit.HOURS,
            intervalUnitCount: 1,
            maximumLines: 10000,
        },
        certificateAutoRenew: {
            enabled: true,
            intervalUnit: TimeUnit.HOURS,
            intervalUnitCount: 1,
        },
        globalBindings: [],
    }
}
