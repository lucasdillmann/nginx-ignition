import I18nContext from "./I18nContext"
import resolveLanguageTag from "./I18nLanguageTagResolver"

export function formatDateTime(value: string | Date): string {
    const { currentLanguage, defaultLanguage, availableLanguages } = I18nContext.get()
    const language = resolveLanguageTag(availableLanguages, currentLanguage, defaultLanguage)

    return new Date(value).toLocaleString(language)
}
