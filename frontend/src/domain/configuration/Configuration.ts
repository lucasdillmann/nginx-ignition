export interface CodeEditor {
    apiKey?: string
}

export interface Version {
    current?: string
    latest?: string
}

export default interface Configuration {
    codeEditor: CodeEditor
    version: Version
}
