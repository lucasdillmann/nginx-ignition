import TrafficStatsGateway from "./TrafficStatsGateway"
import { requireSuccessPayload } from "../../core/apiclient/ApiResponse"
import TrafficStatsResponse from "./model/TrafficStatsResponse"

export default class TrafficStatsService {
    private readonly gateway: TrafficStatsGateway

    constructor() {
        this.gateway = new TrafficStatsGateway()
    }

    async getStats(): Promise<TrafficStatsResponse> {
        return this.gateway.getStats().then(requireSuccessPayload)
    }
}
