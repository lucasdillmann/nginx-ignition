import I18nGateway from "./I18nGateway"
import { requireSuccessPayload } from "../apiclient/ApiResponse"
import I18nContext from "./I18nContext"
import LocalStorageRepository from "../repository/LocalStorageRepository"

export default class I18nService {
    private readonly gateway: I18nGateway
    private readonly repository: LocalStorageRepository<string>

    constructor() {
        this.gateway = new I18nGateway()
        this.repository = new LocalStorageRepository<string>("nginxIgnition.i18n.language")
    }

    async initContext(): Promise<void> {
        const { defaultLanguage, dictionaries } = await this.gateway.getDictionaries().then(requireSuccessPayload)

        I18nContext.replace({
            defaultLanguage,
            dictionaries,
            currentLanguage: this.repository.get(),
        })
    }

    getCustomLanguage(): string | null {
        return this.repository.get()
    }

    setCustomLanguage(languageTag: string | null) {
        if (!languageTag) this.repository.clear()
        else this.repository.set(languageTag)

        I18nContext.replace({
            ...I18nContext.get(),
            currentLanguage: languageTag,
        })
    }
}
