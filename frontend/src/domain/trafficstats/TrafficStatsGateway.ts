import ApiClient from "../../core/apiclient/ApiClient"
import ApiResponse from "../../core/apiclient/ApiResponse"
import TrafficStatsResponse from "./model/TrafficStatsResponse"

export default class TrafficStatsGateway {
    private readonly client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/nginx/traffic-stats")
    }

    async getStats(): Promise<ApiResponse<TrafficStatsResponse>> {
        return this.client.get("")
    }
}
