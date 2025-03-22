import React from "react"
import CodeEditor, { CodeEditorLanguage } from "./CodeEditor"
import { Drawer, Flex, Form, Select } from "antd"

interface CodeEditorModalState {
    language: CodeEditorLanguage
}

export interface CodeEditorModalProps {
    open: boolean
    onClose: () => void
    value: string
    onChange: (value: string) => void
    language: CodeEditorLanguage | CodeEditorLanguage[]
}

export default class CodeEditorModal extends React.Component<CodeEditorModalProps, CodeEditorModalState> {
    constructor(props: CodeEditorModalProps) {
        super(props)

        const { language } = props
        this.state = {
            language: Array.isArray(language) ? language[0] : language,
        }
    }

    private handleLanguageChange(language: CodeEditorLanguage) {
        this.setState({ language })
    }

    private languageName(language: CodeEditorLanguage): string {
        switch (language) {
            case CodeEditorLanguage.JAVASCRIPT:
                return "JavaScript"
            case CodeEditorLanguage.JSON:
                return "JSON"
            case CodeEditorLanguage.HTML:
                return "HTML"
            case CodeEditorLanguage.LUA:
                return "Lua"
            case CodeEditorLanguage.CSS:
                return "CSS"
            case CodeEditorLanguage.YAML:
                return "YAML"
            case CodeEditorLanguage.XML:
                return "XML"
            case CodeEditorLanguage.PLAIN_TEXT:
                return "Plain text"
        }
    }

    private renderLanguageSelector() {
        const { language: availableLanguages } = this.props
        const { language: currentLanguage } = this.state

        if (!Array.isArray(availableLanguages)) {
            return null
        }

        return (
            <Flex justify="right">
                <Form.Item label="Language" layout="horizontal" style={{ width: 250, margin: 0, padding: 0 }} required>
                    <Select
                        onChange={value => this.handleLanguageChange(value)}
                        value={currentLanguage}
                        style={{ width: 150, textAlign: "left", float: "right" }}
                    >
                        {availableLanguages.map(language => (
                            <Select.Option key={language} value={language}>
                                {this.languageName(language)}
                            </Select.Option>
                        ))}
                    </Select>
                </Form.Item>
            </Flex>
        )
    }

    render() {
        const { onClose, value, onChange, open } = this.props
        const { language } = this.state

        if (!open) {
            return undefined
        }

        return (
            <Drawer
                title="Code editor"
                placement="right"
                width="80vw"
                onClose={onClose}
                extra={this.renderLanguageSelector()}
                open
            >
                <CodeEditor value={value} onChange={onChange} language={language} />
            </Drawer>
        )
    }
}
