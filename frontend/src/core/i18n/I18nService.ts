import I18nGateway from "./I18nGateway"
import { requireSuccessPayload } from "../apiclient/ApiResponse"
import I18nContext from "./I18nContext"
import LocalStorageRepository from "../repository/LocalStorageRepository"
import resolveLanguageTag from "./I18nLanguageTagResolver"

export default class I18nService {
    private readonly gateway: I18nGateway
    private readonly repository: LocalStorageRepository<string>

    constructor() {
        this.gateway = new I18nGateway()
        this.repository = new LocalStorageRepository<string>("nginxIgnition.i18n.language")
    }

    async initContext(): Promise<void> {
        const { default: defaultLanguage, available: availableLanguages } = await this.gateway
            .getAvailableLanguages()
            .then(requireSuccessPayload)

        I18nContext.replace({
            ...I18nContext.get(),
            defaultLanguage,
            availableLanguages,
            currentLanguage: this.repository.get(),
        })

        await this.loadDictionary(this.repository.get())
    }

    getCustomLanguage(): string | null {
        return this.repository.get()
    }

    async setCustomLanguage(languageTag: string | null) {
        if (!languageTag) this.repository.clear()
        else this.repository.set(languageTag)

        await this.loadDictionary(languageTag)

        I18nContext.replace({
            ...I18nContext.get(),
            currentLanguage: languageTag,
        })
    }

    private async loadDictionary(targetLanguage: string | null) {
        const context = I18nContext.get()
        const resolvedLanguage = resolveLanguageTag(
            context.availableLanguages,
            targetLanguage ?? context.currentLanguage,
            context.defaultLanguage,
        )

        if (context.loadedDictionaries[resolvedLanguage]) return

        const dictionary = await this.gateway.getDictionary(resolvedLanguage).then(requireSuccessPayload)
        const currentContext = I18nContext.get()

        I18nContext.replace({
            ...currentContext,
            loadedDictionaries: {
                ...currentContext.loadedDictionaries,
                [resolvedLanguage]: dictionary,
            },
        })
    }
}
