import React from "react"
import { Form, Modal, Select, Switch } from "antd"
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

    private cancel() {
        const { onClose } = this.props
        onClose()
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
            this.setState({
                config: this.service.currentConfig(),
            })
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
            >
                <Form
                    {...FormStyle}
                    initialValues={config}
                    onValuesChange={(_, values) => this.handleChange(values)}
                    style={{
                        marginTop: 30,
                        marginBottom: 50,
                    }}
                >
                    <Form.Item
                        label={<I18n id={MessageKey.FrontendComponentsDatatableOptionsMode} />}
                        name="mode"
                        required
                    >
                        <Select>
                            <Select.Option value={DataTablePersistentStateMode.GLOBAL}>
                                <I18n id={MessageKey.FrontendComponentsDatatableOptionsModeGlobal} />
                            </Select.Option>
                            <Select.Option value={DataTablePersistentStateMode.BY_TABLE}>
                                <I18n id={MessageKey.FrontendComponentsDatatableOptionsModeByTable} />
                            </Select.Option>
                            <Select.Option value={DataTablePersistentStateMode.FIXED}>
                                <I18n id={MessageKey.FrontendComponentsDatatableOptionsModeFixed} />
                            </Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item
                        label={<I18n id={MessageKey.FrontendComponentsDatatableOptionsPageSize} />}
                        name="pageSize"
                        required
                    >
                        <Select>
                            {DATA_TABLE_PAGE_SIZES.map(size => (
                                <Select.Option key={size} value={size}>
                                    {size}
                                </Select.Option>
                            ))}
                        </Select>
                    </Form.Item>
                    <Form.Item
                        label={<I18n id={MessageKey.FrontendComponentsDatatableOptionsPersistPageNumber} />}
                        name="persistPageNumber"
                        valuePropName="checked"
                        help={<I18n id={MessageKey.FrontendComponentsDatatableOptionsPersistPageNumberHelp} />}
                        required
                    >
                        <Switch />
                    </Form.Item>
                    <Form.Item
                        label={<I18n id={MessageKey.FrontendComponentsDatatableOptionsPersistSearchTerms} />}
                        name="persistSearchTerms"
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
