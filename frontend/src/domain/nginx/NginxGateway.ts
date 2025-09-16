import ApiClient from "../../core/apiclient/ApiClient"
import ApiResponse from "../../core/apiclient/ApiResponse"
import { NginxStatusResponse } from "./model/NginxStatusResponse"

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

    async getLogs(lines: number): Promise<ApiResponse<string[]>> {
        return this.client.get("/logs", undefined, { lines })
    }

    async configFiles(configPath: string, logPath: string): Promise<ApiResponse<any>> {
        return this.client.get("/config", undefined, { configPath, logPath }, true)
    }
}
