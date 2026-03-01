import ApiClient from "../../core/apiclient/ApiClient"
import ApiResponse from "../../core/apiclient/ApiResponse"
import I18nAvailableLanguages from "./model/I18nAvailableLanguages"
import I18nDictionary from "./model/I18nDictionary"

export default class I18nGateway {
    private readonly client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/i18n")
    }

    async getAvailableLanguages(): Promise<ApiResponse<I18nAvailableLanguages>> {
        return this.client.get()
    }

    async getDictionary(language: string): Promise<ApiResponse<I18nDictionary>> {
        return this.client.get(`/${language}`)
    }
}
