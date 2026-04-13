import React from "react"
import { Form, FormInstance, Modal, Select, Switch } from "antd"
import DataTablePersistentStateConfig from "./model/DataTablePersistentStateConfig"
import DataTableService from "./DataTableService"
import { DataTablePersistentStateMode } from "./model/DataTablePersistentStateMode"
import { DATA_TABLE_PAGE_SIZES } from "./model/DataTablePageSize"
import { I18n } from "../../i18n/I18n"
import MessageKey from "../../i18n/model/MessageKey.generated"
import FormLayout from "../form/FormLayout"

const FormStyle = {
    ...FormLayout.FormDefaults,
    labelCol: { span: 8 },
    wrapperCol: { span: 16 },
}

const PAGINATION_MODE_OPTIONS_DATA = [
    { value: DataTablePersistentStateMode.GLOBAL, messageKey: MessageKey.FrontendComponentsDatatableOptionsModeGlobal },
    {
        value: DataTablePersistentStateMode.BY_TABLE,
        messageKey: MessageKey.FrontendComponentsDatatableOptionsModeByTable,
    },
    { value: DataTablePersistentStateMode.FIXED, messageKey: MessageKey.FrontendComponentsDatatableOptionsModeFixed },
]

const PAGE_SIZE_OPTIONS = DATA_TABLE_PAGE_SIZES.map(size => ({
    value: size,
    label: size,
}))

export interface DataTableOptionsProps {
    id: string
    open: boolean
    onClose: () => void
}

interface DataTableOptionsState {
    config: DataTablePersistentStateConfig
}

export default class DataTableOptions extends React.Component<DataTableOptionsProps, DataTableOptionsState> {
    private readonly service: DataTableService
    private readonly formRef = React.createRef<FormInstance>()

    constructor(props: DataTableOptionsProps) {
        super(props)
        this.service = new DataTableService()

        this.state = {
            config: this.service.currentConfig(),
        }
    }

    private save() {
        const { config } = this.state
        this.service.updateConfig(config)

        const { onClose } = this.props
        onClose()
    }

    private refreshForm() {
        const { config } = this.state
        this.formRef.current?.setFieldsValue(config)
    }

    private cancel() {
        const { onClose } = this.props

        this.setState(
            {
                config: this.service.currentConfig(),
            },
            () => {
                this.refreshForm()
                onClose()
            },
        )
    }

    private handleChange(newConfig: Partial<DataTablePersistentStateConfig>) {
        const { config } = this.state
        this.setState({
            config: {
                ...config,
                ...newConfig,
            },
        })
    }

    componentDidUpdate(prevProps: Readonly<DataTableOptionsProps>) {
        const { open } = this.props

        if (prevProps.open !== open && open) {
            this.setState(
                {
                    config: this.service.currentConfig(),
                },
                () => this.refreshForm(),
            )
        }
    }

    render() {
        const { open } = this.props
        const { config } = this.state

        return (
            <Modal
                title={<I18n id={MessageKey.FrontendComponentsDatatableOptionsTitle} />}
                open={open}
                onCancel={() => this.cancel()}
                onOk={() => this.save()}
                width={700}
                forceRender
            >
                <Form
                    {...FormStyle}
                    ref={this.formRef}
                    initialValues={config}
                    onValuesChange={(_, values) => this.handleChange(values)}
                    style={{
                        marginTop: 30,
                        marginBottom: 50,
                    }}
                >
                    <Form.Item
                        label={<I18n id={MessageKey.FrontendComponentsDatatableOptionsMode} />}
                        name="paginationMode"
                        required
                    >
                        <Select
                            options={PAGINATION_MODE_OPTIONS_DATA.map(item => ({
                                value: item.value,
                                label: <I18n id={item.messageKey} />,
                            }))}
                        />
                    </Form.Item>
                    <Form.Item
                        label={<I18n id={MessageKey.FrontendComponentsDatatableOptionsPageSize} />}
                        name="defaultPageSize"
                        required
                    >
                        <Select options={PAGE_SIZE_OPTIONS} />
                    </Form.Item>
                    <Form.Item
                        label={<I18n id={MessageKey.FrontendComponentsDatatableOptionsPersistPageNumber} />}
                        name="rememberPageNumber"
                        valuePropName="checked"
                        help={<I18n id={MessageKey.FrontendComponentsDatatableOptionsPersistPageNumberHelp} />}
                        required
                    >
                        <Switch />
                    </Form.Item>
                    <Form.Item
                        label={<I18n id={MessageKey.FrontendComponentsDatatableOptionsPersistSearchTerms} />}
                        name="rememberSearchTerms"
                        valuePropName="checked"
                        help={<I18n id={MessageKey.FrontendComponentsDatatableOptionsPersistSearchTermsHelp} />}
                        required
                    >
                        <Switch />
                    </Form.Item>
                </Form>
            </Modal>
        )
    }
}
