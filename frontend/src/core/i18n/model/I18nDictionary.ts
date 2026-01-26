import MessageKey from "./MessageKey.generated"

export default interface I18nDictionary {
    languageTag: string
    messages: Record<MessageKey, string>
}
