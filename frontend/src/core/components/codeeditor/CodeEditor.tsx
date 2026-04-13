import React from "react"
import CodeMirror from "@uiw/react-codemirror"
import type { Extension } from "@codemirror/state"
import { StreamLanguage, syntaxHighlighting, defaultHighlightStyle } from "@codemirror/language"
import { javascript } from "@codemirror/lang-javascript"
import { json } from "@codemirror/lang-json"
import { html } from "@codemirror/lang-html"
import { css } from "@codemirror/lang-css"
import { xml } from "@codemirror/lang-xml"
import { yaml } from "@codemirror/lang-yaml"
import { lua } from "@codemirror/legacy-modes/mode/lua"

export enum CodeEditorLanguage {
    JAVASCRIPT = "javascript",
    JSON = "json",
    HTML = "html",
    LUA = "lua",
    CSS = "css",
    YAML = "yaml",
    XML = "xml",
    PLAIN_TEXT = "plaintext",
}

export interface CodeEditorProps {
    value: string
    onChange: (value: string) => void
    language: CodeEditorLanguage
}

function languageExtensions(language: CodeEditorLanguage): Extension[] {
    switch (language) {
        case CodeEditorLanguage.JAVASCRIPT:
            return [javascript()]
        case CodeEditorLanguage.JSON:
            return [json()]
        case CodeEditorLanguage.HTML:
            return [html()]
        case CodeEditorLanguage.LUA:
            return [StreamLanguage.define(lua), syntaxHighlighting(defaultHighlightStyle)]
        case CodeEditorLanguage.CSS:
            return [css()]
        case CodeEditorLanguage.YAML:
            return [yaml()]
        case CodeEditorLanguage.XML:
            return [xml()]
        case CodeEditorLanguage.PLAIN_TEXT:
            return []
    }
}

export default class CodeEditor extends React.Component<CodeEditorProps> {
    private handleOnChange(value: string) {
        const { onChange } = this.props
        onChange(value)
    }

    render() {
        const { language, value } = this.props

        return (
            <CodeMirror
                key={language}
                value={value}
                height="100%"
                theme="dark"
                extensions={languageExtensions(language)}
                onChange={value => this.handleOnChange(value)}
                onCreateEditor={view => view.focus()}
            />
        )
    }
}
