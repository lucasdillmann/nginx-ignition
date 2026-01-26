import React from "react"
import CodeEditor, { CodeEditorLanguage } from "./CodeEditor"
import { Drawer, Flex, Form, Select } from "antd"
import { i18n, I18n } from "../../i18n/I18n"
import MessageKey from "../../i18n/model/MessageKey.generated"

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

    private languageName(language: CodeEditorLanguage): MessageKey {
        switch (language) {
            case CodeEditorLanguage.JAVASCRIPT:
                return MessageKey.FrontendComponentsCodeeditorLanguageJavascript
            case CodeEditorLanguage.JSON:
                return MessageKey.FrontendComponentsCodeeditorLanguageJson
            case CodeEditorLanguage.HTML:
                return MessageKey.FrontendComponentsCodeeditorLanguageHtml
            case CodeEditorLanguage.LUA:
                return MessageKey.FrontendComponentsCodeeditorLanguageLua
            case CodeEditorLanguage.CSS:
                return MessageKey.FrontendComponentsCodeeditorLanguageCss
            case CodeEditorLanguage.YAML:
                return MessageKey.FrontendComponentsCodeeditorLanguageYaml
            case CodeEditorLanguage.XML:
                return MessageKey.FrontendComponentsCodeeditorLanguageXml
            case CodeEditorLanguage.PLAIN_TEXT:
                return MessageKey.FrontendComponentsCodeeditorLanguagePlainText
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
                <Form.Item
                    label={<I18n id={MessageKey.CommonLanguage} />}
                    layout="horizontal"
                    style={{ width: 250, margin: 0, padding: 0 }}
                    required
                >
                    <Select
                        onChange={value => this.handleLanguageChange(value)}
                        value={currentLanguage}
                        style={{ width: 150, textAlign: "left", float: "right" }}
                    >
                        {availableLanguages.map(language => (
                            <Select.Option key={language} value={language}>
                                <I18n id={this.languageName(language)} />
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
                title={i18n(MessageKey.FrontendComponentsCodeeditorTitle)}
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
