import NginxGateway from "../nginx/NginxGateway"
import BackupGateway from "../backup/BackupGateway"
import { requireSuccessRawResponse } from "../../core/apiclient/ApiResponse"

export default class ExportService {
    private readonly nginxGateway: NginxGateway
    private readonly backupGateway: BackupGateway

    constructor() {
        this.nginxGateway = new NginxGateway()
        this.backupGateway = new BackupGateway()
    }

    async downloadNginxConfigurationFiles(basePath: string, configPath: string, logPath: string): Promise<void> {
        return this.nginxGateway
            .configFiles(basePath, configPath, logPath)
            .then(requireSuccessRawResponse)
            .then(response => response.raw.blob())
            .then(blob => this.sendBlob(blob, "nginx-config.zip"))
    }

    async downloadDatabaseBackup(): Promise<void> {
        return this.backupGateway
            .download()
            .then(requireSuccessRawResponse)
            .then(response => Promise.all([response.raw.blob(), response.raw]))
            .then(([blob, response]) => ({
                blob,
                fileName: this.getFileName(response),
            }))
            .then(data => this.sendBlob(data.blob, data.fileName))
    }

    private getFileName(response: Response): string {
        const fallbackName = "backup.bin"

        const contentDisposition = response.headers.get("content-disposition")
        if (!contentDisposition) return fallbackName

        const fileNameMatch = /filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/.exec(contentDisposition)
        if (fileNameMatch?.[1]) {
            return fileNameMatch[1].replace(/['"]/g, "")
        }

        return fallbackName
    }

    private async sendBlob(blob: Blob, fileName: string): Promise<void> {
        const downloadUrl = window.URL.createObjectURL(blob)
        const link = document.createElement("a")

        link.href = downloadUrl
        link.download = fileName

        document.body.appendChild(link)
        link.click()

        document.body.removeChild(link)
        window.URL.revokeObjectURL(downloadUrl)
    }
}
