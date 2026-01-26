import ApiClient from "../../core/apiclient/ApiClient"
import ApiResponse from "../../core/apiclient/ApiResponse"
import I18nDictionaries from "./model/I18nDictionaries"

export default class I18nGateway {
    private readonly client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/i18n")
    }

    async getDictionaries(): Promise<ApiResponse<I18nDictionaries>> {
        return this.client.get()
    }
}
