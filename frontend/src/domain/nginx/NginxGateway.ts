import ApiClient from "../../core/apiclient/ApiClient"
import ApiResponse from "../../core/apiclient/ApiResponse"
import { NginxStatusResponse } from "./model/NginxStatusResponse"
import NginxMetadata from "./model/NginxMetadata"

export default class NginxGateway {
    private readonly client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/nginx")
    }

    async start(): Promise<ApiResponse<void>> {
        return this.client.post("/start")
    }

    async stop(): Promise<ApiResponse<void>> {
        return this.client.post("/stop")
    }

    async reload(): Promise<ApiResponse<void>> {
        return this.client.post("/reload")
    }

    async getStatus(): Promise<ApiResponse<NginxStatusResponse>> {
        return this.client.get("/status")
    }

    async getMetadata(): Promise<ApiResponse<NginxMetadata>> {
        return this.client.get("/metadata")
    }

    async getLogs(lines: number): Promise<ApiResponse<string[]>> {
        return this.client.get("/logs", undefined, { lines })
    }

    async configFiles(
        basePath: string,
        configPath: string,
        logPath: string,
        cachePath: string,
        tempPath: string,
    ): Promise<ApiResponse<any>> {
        return this.client.get("/config", undefined, { basePath, configPath, logPath, cachePath, tempPath }, true)
    }
}
