import NginxGateway from "../nginx/NginxGateway"

export default class ExportService {
    private readonly nginxGateway: NginxGateway

    constructor() {
        this.nginxGateway = new NginxGateway()
    }

    async downloadNginxConfigurationFiles(): Promise<void> {
        return this.nginxGateway
            .configFiles()
            .then(response => response.raw.blob())
            .then(blob => {
                const downloadUrl = window.URL.createObjectURL(blob)
                const link = document.createElement("a")

                link.href = downloadUrl
                link.download = "nginx-config.zip"

                document.body.appendChild(link)
                link.click()

                document.body.removeChild(link)
                window.URL.revokeObjectURL(downloadUrl)
            })
    }

    async downloadDatabaseBackup(): Promise<void> {
        return Promise.reject(new Error("Not implemented yet"))
    }
}
