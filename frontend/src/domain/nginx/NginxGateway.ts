import ApiClient from "../../core/apiclient/ApiClient"
import ApiResponse from "../../core/apiclient/ApiResponse"
import { NginxStatusResponse } from "./model/NginxStatusResponse"

export default class NginxGateway {
    private client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/nginx")
    }

    async start(): Promise<ApiResponse<undefined>> {
        return this.client.post("/start")
    }

    async stop(): Promise<ApiResponse<undefined>> {
        return this.client.post("/stop")
    }

    async reload(): Promise<ApiResponse<undefined>> {
        return this.client.post("/reload")
    }

    async getStatus(): Promise<ApiResponse<NginxStatusResponse>> {
        return this.client.get("/status")
    }

    async getLogs(lines: number): Promise<ApiResponse<string[]>> {
        return this.client.get("/logs", undefined, { lines })
    }
}
