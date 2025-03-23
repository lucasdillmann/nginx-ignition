import NginxGateway from "./NginxGateway"
import { requireSuccessPayload, requireSuccessResponse } from "../../core/apiclient/ApiResponse"
import NginxEventDispatcher from "./listener/NginxEventDispatcher"
import { NginxOperation } from "./listener/NginxEventListener"

export default class NginxService {
    private readonly gateway: NginxGateway

    constructor() {
        this.gateway = new NginxGateway()
    }

    async isRunning(): Promise<boolean> {
        return this.gateway
            .getStatus()
            .then(requireSuccessPayload)
            .then(response => response.running)
    }

    async start(): Promise<void> {
        return this.gateway
            .start()
            .then(requireSuccessResponse)
            .then(response => {
                NginxEventDispatcher.notify(NginxOperation.START)
                return response
            })
    }

    async stop(): Promise<void> {
        return this.gateway
            .stop()
            .then(requireSuccessResponse)
            .then(response => {
                NginxEventDispatcher.notify(NginxOperation.STOP)
                return response
            })
    }

    async reloadConfiguration(): Promise<void> {
        return this.gateway
            .reload()
            .then(requireSuccessResponse)
            .then(response => {
                NginxEventDispatcher.notify(NginxOperation.RELOAD)
                return response
            })
    }

    async logs(lines: number): Promise<string[]> {
        return this.gateway.getLogs(lines).then(requireSuccessPayload)
    }
}
