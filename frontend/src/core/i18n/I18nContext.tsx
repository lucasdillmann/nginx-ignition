import ContextHolder from "../context/ContextHolder"
import I18nDictionary from "./model/I18nDictionary"

export interface I18nContext {
    defaultLanguage: string
    currentLanguage: string | null
    dictionaries: I18nDictionary[]
}

export default new ContextHolder<I18nContext>({
    defaultLanguage: "",
    currentLanguage: null,
    dictionaries: [],
})
