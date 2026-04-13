import React from "react"
import Editor from "@monaco-editor/react"
import { editor } from "monaco-editor"

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

export default class CodeEditor extends React.Component<CodeEditorProps> {
    private handleOnChange(value: string | undefined) {
        const { onChange } = this.props
        onChange(value ?? "")
    }

    private handleMount(editor: editor.IStandaloneCodeEditor) {
        editor.focus()
    }

    render() {
        const { language, value } = this.props

        return (
            <Editor
                language={language}
                theme="vs-dark"
                onChange={value => this.handleOnChange(value)}
                value={value}
                onMount={editor => this.handleMount(editor)}
                height="100%"
                width="100%"
                options={{ automaticLayout: true }}
            />
        )
    }
}
