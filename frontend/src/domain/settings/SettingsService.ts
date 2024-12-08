import { requireSuccessPayload, requireSuccessResponse } from "../../core/apiclient/ApiResponse"
import SettingsDto from "./model/SettingsDto"
import SettingsGateway from "./SettingsGateway"

export default class SettingsService {
    private readonly gateway: SettingsGateway

    constructor() {
        this.gateway = new SettingsGateway()
    }

    async get(): Promise<SettingsDto> {
        return this.gateway.get().then(requireSuccessPayload)
    }

    async save(settings: SettingsDto): Promise<void> {
        return this.gateway.put(settings).then(requireSuccessResponse)
    }
}
