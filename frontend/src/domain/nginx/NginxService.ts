import NginxGateway from "./NginxGateway"
import { requireSuccessPayload, requireSuccessResponse } from "../../core/apiclient/ApiResponse"
import NginxEventDispatcher from "./listener/NginxEventDispatcher"
import { NginxOperation } from "./listener/NginxEventListener"
import NginxMetadata from "./model/NginxMetadata"
import LogLine from "../logs/model/LogLine"

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

    async getMetadata(): Promise<NginxMetadata> {
        return this.gateway.getMetadata().then(requireSuccessPayload)
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

    async logs(lines: number, surroundingLines: number, searchTerms?: string): Promise<LogLine[]> {
        return this.gateway.getLogs(lines, surroundingLines, searchTerms).then(requireSuccessPayload)
    }
}
