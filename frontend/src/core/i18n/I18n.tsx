import MessageKey from "./model/MessageKey.generated"
import React from "react"
import I18nDictionary from "./model/I18nDictionary"
import I18nContext from "./I18nContext"

interface I18nState {
    value: string
}

export interface MessageKeyWithParams {
    id: MessageKey
    params: Record<string, any>
}

export type I18nMessage = MessageKey | MessageKeyWithParams

export interface I18nProps {
    id: I18nMessage
    params?: Record<string, any>
    fallback?: string
}

export class I18n extends React.Component<I18nProps, I18nState> {
    constructor(props: I18nProps) {
        super(props)

        this.state = {
            value: this.resolveMessage(),
        }
    }

    private resolveMessage(): string {
        const { id, params, fallback } = this.props
        return params && typeof id === "string" ? i18n({ id, params }, fallback) : i18n(id, fallback)
    }

    private handleContextUpdate() {
        this.setState({
            value: this.resolveMessage(),
        })
    }

    componentDidMount() {
        I18nContext.register(this.handleContextUpdate.bind(this))
    }

    componentWillUnmount() {
        I18nContext.deregister(this.handleContextUpdate.bind(this))
    }

    render() {
        const { value } = this.state
        return value
    }
}

export function raw(message: string): I18nMessage {
    return { id: MessageKey.CommonRaw, params: { message } }
}

export function i18n(input: I18nMessage, fallback?: string): string {
    let id: MessageKey
    let params: Record<string, any> = {}

    if (typeof input === "object") ({ id, params } = input)
    else id = input

    const dictionary = resolveDictionary()
    const template = dictionary.messages[id]
    if (!template && template !== "") return fallback ?? id

    return template.replace(/\${(.*?)}/g, (match, varName) => {
        const value = params[varName]
        return value !== undefined ? String(value) : match
    })
}

function resolveDictionary(): I18nDictionary {
    const { currentLanguage, defaultLanguage, dictionaries } = I18nContext.get()
    const targetLanguage = currentLanguage ?? defaultLanguage

    for (const dictionary of dictionaries) {
        if (dictionary.languageTag === targetLanguage) {
            return dictionary
        }
    }

    if (targetLanguage.includes("-")) {
        const baseLanguage = targetLanguage.split("-")[0]
        for (const dictionary of dictionaries) {
            if (dictionary.languageTag === baseLanguage) {
                return dictionary
            }
        }
    }

    for (const dictionary of dictionaries) {
        const baseLanguage = defaultLanguage.split("-")[0]
        if (dictionary.languageTag === defaultLanguage || dictionary.languageTag === baseLanguage) {
            return dictionary
        }
    }

    // @ts-expect-error fallback dictionary
    return { languageTag: targetLanguage, messages: {} }
}
