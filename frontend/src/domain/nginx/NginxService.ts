import NginxGateway from "./NginxGateway";
import {requireSuccessPayload, requireSuccessResponse} from "../../core/apiclient/ApiResponse";

export default class NginxService {
    private gateway: NginxGateway

    constructor() {
        this.gateway = new NginxGateway()
    }

    async isRunning(): Promise<boolean> {
        return this.gateway
            .getStatus()
            .then(requireSuccessPayload)
            .then(response => response.running)
    }

    async start(): Promise<undefined> {
        return this.gateway.start().then(requireSuccessResponse)
    }

    async stop(): Promise<undefined> {
        return this.gateway.stop().then(requireSuccessResponse)
    }

    async reloadConfiguration(): Promise<undefined> {
        return this.gateway.reload().then(requireSuccessResponse)
    }
}
