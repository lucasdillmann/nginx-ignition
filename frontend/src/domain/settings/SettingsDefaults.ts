import SettingsFormValues from "./model/SettingsFormValues"
import { LogLevel, TimeUnit } from "./model/SettingsDto"

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
                clientBody: 60,
            },
            buffers: {
                clientBodyKb: 16,
                clientHeaderKb: 1,
                largeClientHeader: {
                    amount: 4,
                    sizeKb: 8,
                },
                output: {
                    amount: 4,
                    sizeKb: 32,
                },
            },
            defaultContentType: "application/octet-stream",
            gzipEnabled: true,
            maximumBodySizeMb: 1024,
            sendfileEnabled: true,
            tcpNoDelayEnabled: true,
            serverTokensEnabled: false,
            workerConnections: 1024,
            workerProcesses: 2,
            runtimeUser: "root",
            custom: null,
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
