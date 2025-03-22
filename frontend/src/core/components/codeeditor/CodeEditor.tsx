import React from "react"
import { CodeiumEditor } from "@codeium/react-code-editor"
import AppContext from "../context/AppContext"

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

export default class CodeEditor extends React.Component<CodeEditorProps, any> {
    private handleOnChange(value: string | undefined) {
        const { onChange } = this.props

        const newValue = value ?? ""
        onChange(newValue)
    }

    render() {
        const { language, value } = this.props
        const { configuration } = AppContext.get()

        return (
            <CodeiumEditor
                language={language}
                theme="vs-dark"
                onChange={value => this.handleOnChange(value)}
                value={value}
                apiKey={configuration.codeEditor.apiKey}
                onMount={editor => editor.focus()}
                height="100%"
                width="100%"
            />
        )
    }
}
