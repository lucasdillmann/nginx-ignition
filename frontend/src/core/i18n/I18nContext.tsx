import ContextHolder from "../context/ContextHolder"
import I18nDictionary from "./model/I18nDictionary"

export interface I18nContext {
    defaultLanguage: string
    currentLanguage: string | null
    availableLanguages: string[]
    loadedDictionaries: Record<string, I18nDictionary>
}

export default new ContextHolder<I18nContext>({
    defaultLanguage: window.navigator.language,
    currentLanguage: window.navigator.languages.join(","),
    availableLanguages: [],
    loadedDictionaries: {},
})
