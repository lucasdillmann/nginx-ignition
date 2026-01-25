import React, { CSSProperties } from "react"
import I18nService from "./I18nService"
import I18nContext from "./I18nContext"
import { I18n } from "./I18n"
import MessageKey from "./model/MessageKey.generated"
import { Form, Modal, Select } from "antd"
import { TranslationOutlined } from "@ant-design/icons"
import FormLayout from "../components/form/FormLayout"

interface AvailableLanguage {
    label: React.ReactNode
    value: string | null
}

interface I18nLanguagePickerState {
    open: boolean
    current: string | null
    available: AvailableLanguage[]
}

export interface I18nLanguagePickerProps {
    style?: CSSProperties
}

export default class I18nLanguagePicker extends React.Component<I18nLanguagePickerProps, I18nLanguagePickerState> {
    private readonly service: I18nService

    constructor(props: I18nLanguagePickerProps) {
        super(props)
        this.service = new I18nService()

        const available: AvailableLanguage[] = I18nContext.get().dictionaries.map(d => ({
            label: <I18n id={`frontend/i18n/lang-name-${d.languageTag}` as MessageKey} />,
            value: d.languageTag,
        }))

        available.unshift({
            label: <I18n id={MessageKey.FrontendI18nAuto} />,
            value: null,
        })

        this.state = {
            available,
            open: false,
            current: this.service.getCustomLanguage(),
        }
    }

    private handleLanguageChange(languageTag: string | null) {
        this.service.setCustomLanguage(languageTag)
        this.setState({ current: languageTag })
    }

    private changeVisibility(open: boolean) {
        this.setState({ open })
    }

    render() {
        const { style } = this.props
        const { available, current, open } = this.state
        const selected = available.find(l => l.value === current)

        // TODO: Add beta badge and additional instructions (link to report translation issues and alike)

        return (
            <>
                <TranslationOutlined onClick={() => this.changeVisibility(true)} style={style} />

                <Modal
                    title={<I18n id={MessageKey.FrontendI18nChangeLanguage} />}
                    open={open}
                    afterClose={() => this.changeVisibility(false)}
                    onCancel={() => this.changeVisibility(false)}
                    footer={null}
                >
                    <br />
                    <Form {...FormLayout.FormDefaults} layout="horizontal">
                        <Form.Item
                            label={<I18n id={MessageKey.CommonLanguage} />}
                            help={<I18n id={MessageKey.FrontendI18nChangeLanguageHelp} />}
                            required
                        >
                            <Select
                                options={available}
                                value={selected}
                                onChange={(_, newValue) =>
                                    this.handleLanguageChange((newValue as AvailableLanguage).value)
                                }
                                style={{ width: "100%" }}
                                autoFocus
                            />
                        </Form.Item>
                    </Form>
                </Modal>
            </>
        )
    }
}
