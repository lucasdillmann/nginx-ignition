import React from "react"
import { CodeiumEditor } from "@codeium/react-code-editor"
import AppContext from "../context/AppContext"

export enum CodeEditorLanguage {
    JAVASCRIPT = "JAVASCRIPT",
    JSON = "JSON",
    HTML = "HTML",
    LUA = "LUA",
    CSS = "CSS",
    YAML = "YAML",
    XML = "XML",
    PLAIN_TEXT = "PLAINTEXT",
    GENERIC = "UNSPECIFIED",
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
            />
        )
    }
}
