import ApiClient from "../../core/apiclient/ApiClient"
import ApiResponse from "../../core/apiclient/ApiResponse"
import Configuration from "./Configuration"

export default class ConfigurationGateway {
    private readonly client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/frontend/configuration")
    }

    public get(): Promise<ApiResponse<Configuration>> {
        return this.client.get()
    }
}
