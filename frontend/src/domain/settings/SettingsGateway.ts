import ApiClient from "../../core/apiclient/ApiClient"
import ApiResponse from "../../core/apiclient/ApiResponse"
import SettingsDto from "./model/SettingsDto"

export default class SettingsGateway {
    private readonly client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/settings")
    }

    async get(): Promise<ApiResponse<SettingsDto>> {
        return this.client.get<SettingsDto>()
    }

    async put(settings: SettingsDto): Promise<ApiResponse<void>> {
        return this.client.put(undefined, settings)
    }
}
