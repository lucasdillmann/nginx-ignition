import NginxGateway from "./NginxGateway"
import { requireSuccessPayload } from "../../core/apiclient/ApiResponse"
import NginxEventDispatcher from "./listener/NginxEventDispatcher"
import { NginxOperation } from "./listener/NginxEventListener"
import { NginxActionResponse } from "./model/NginxActionResponse"

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

    async start(): Promise<NginxActionResponse> {
        return this.gateway
            .start()
            .then(requireSuccessPayload)
            .then(response => {
                NginxEventDispatcher.notify(NginxOperation.START)
                return response
            })
    }

    async stop(): Promise<NginxActionResponse> {
        return this.gateway
            .stop()
            .then(requireSuccessPayload)
            .then(response => {
                NginxEventDispatcher.notify(NginxOperation.STOP)
                return response
            })
    }

    async reloadConfiguration(): Promise<NginxActionResponse> {
        return this.gateway
            .reload()
            .then(requireSuccessPayload)
            .then(response => {
                NginxEventDispatcher.notify(NginxOperation.RELOAD)
                return response
            })
    }

    async logs(lines: number): Promise<string[]> {
        return this.gateway.getLogs(lines).then(requireSuccessPayload)
    }
}
