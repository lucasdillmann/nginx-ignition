import React from "react"
import ValidationResult from "../../../core/validation/ValidationResult"
import { Button, Flex, Form, FormListFieldData, FormListOperation, Input, InputNumber, Modal, Select } from "antd"
import If from "../../../core/components/flowcontrol/If"
import { CloseOutlined, PlusOutlined, ArrowUpOutlined, ArrowDownOutlined, SettingOutlined } from "@ant-design/icons"
import { HostFormRoute } from "../model/HostFormValues"
import { HostRouteType } from "../model/HostRequest"
import FormLayout from "../../../core/components/form/FormLayout"
import "./HostRoutes.css"
import TextArea from "antd/es/input/TextArea"

const ACTION_ICON_STYLE = {
    marginLeft: 15,
    alignItems: "start",
    marginTop: 37,
}

const DISABLED_ACTION_ICON_STYLE = {
    ...ACTION_ICON_STYLE,
    color: "#a8a8a8",
}

const ENABLED_ACTION_ICON_STYLE = {
    ...ACTION_ICON_STYLE,
    color: "#000000",
}

const DEFAULT_VALUES: HostFormRoute = {
    priority: 0,
    type: HostRouteType.PROXY,
    sourcePath: "/",
    targetUri: "",
}

export interface HostRoutesProps {
    routes: HostFormRoute[]
    validationResult: ValidationResult
}

interface HostRoutesState {
    customSettingsOpenModalIndex?: number
}

export default class HostRoutes extends React.Component<HostRoutesProps, HostRoutesState> {
    constructor(props: HostRoutesProps) {
        super(props)
        this.state = {}
    }

    private renderProxyRoute(field: FormListFieldData, index: number): React.ReactNode {
        const { validationResult } = this.props
        const { name } = field

        return (
            <Form.Item
                {...FormLayout.ExpandedLabeledItem}
                className="host-form-route-target-uri"
                layout="vertical"
                name={[name, "targetUri"]}
                validateStatus={validationResult.getStatus(`routes[${index}].targetUri`)}
                help={validationResult.getMessage(`routes[${index}].targetUri`)}
                label="Destination URL"
                required
            >
                <Input />
            </Form.Item>
        )
    }

    private renderRedirectRoute(field: FormListFieldData, index: number): React.ReactNode {
        const { validationResult } = this.props
        const { name } = field

        return (
            <>
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-route-status-code"
                    layout="vertical"
                    name={[name, "statusCode"]}
                    validateStatus={validationResult.getStatus(`routes[${index}].statusCode`)}
                    help={validationResult.getMessage(`routes[${index}].statusCode`)}
                    label="Status code"
                    required
                >
                    <InputNumber min={300} max={399} />
                </Form.Item>
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-route-target-uri"
                    layout="vertical"
                    name={[name, "targetUri"]}
                    validateStatus={validationResult.getStatus(`routes[${index}].targetUri`)}
                    help={validationResult.getMessage(`routes[${index}].targetUri`)}
                    label="Destination URL"
                    required
                >
                    <Input />
                </Form.Item>
            </>
        )
    }

    private renderStaticResponseRoute(field: FormListFieldData, index: number): React.ReactNode {
        const { validationResult } = this.props
        const { name } = field

        return (
            <>
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-route-status-code"
                    layout="vertical"
                    name={[name, "response", "statusCode"]}
                    validateStatus={validationResult.getStatus(`routes[${index}].response.statusCode`)}
                    help={validationResult.getMessage(`routes[${index}].response.statusCode`)}
                    label="Status code"
                    required
                >
                    <InputNumber min={100} max={599} />
                </Form.Item>
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-route-response-headers"
                    layout="vertical"
                    name={[name, "response", "headers"]}
                    validateStatus={validationResult.getStatus(`routes[${index}].response.headers`)}
                    help={
                        validationResult.getMessage(`routes[${index}].response.headers`) ??
                        "One per line, as [key]: [value]"
                    }
                    label="Headers"
                    required
                >
                    <TextArea rows={3} />
                </Form.Item>
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-route-response-payload"
                    layout="vertical"
                    name={[name, "response", "payload"]}
                    validateStatus={validationResult.getStatus(`routes[${index}].response.payload`)}
                    help={validationResult.getMessage(`routes[${index}].response.payload`)}
                    label="Body / Payload"
                    required
                >
                    <TextArea rows={3} />
                </Form.Item>
            </>
        )
    }

    private moveRoute(operations: FormListOperation, index: number, offset: number) {
        const { routes } = this.props

        if (index === 0 && offset < 0) return
        if (index === routes.length && offset > 0) return

        const currentPosition = routes[index]
        const newPosition = routes[index + offset]

        currentPosition.priority = currentPosition.priority + offset
        newPosition.priority = newPosition.priority - offset

        const indexToRemove = offset > 0 ? index : index + offset
        operations.remove(indexToRemove)
        operations.remove(indexToRemove)
        operations.add(currentPosition)
        operations.add(newPosition)
    }

    private openCustomSettingsModal(index: number) {
        this.setState({ customSettingsOpenModalIndex: index })
    }

    private closeCustomSettingsModal() {
        this.setState({ customSettingsOpenModalIndex: undefined })
    }

    private renderRoute(field: FormListFieldData, operations: FormListOperation, index: number) {
        const { validationResult, routes } = this.props
        const { customSettingsOpenModalIndex } = this.state
        const { name } = field
        const type = routes[index].type

        return (
            <Flex className="host-form-route-container">
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-route-type"
                    layout="vertical"
                    name={[name, "type"]}
                    validateStatus={validationResult.getStatus(`routes[${index}].type`)}
                    help={validationResult.getMessage(`routes[${index}].type`)}
                    label="Type"
                    required
                >
                    <Select>
                        <Select.Option value={HostRouteType.PROXY}>Proxy</Select.Option>
                        <Select.Option value={HostRouteType.REDIRECT}>Redirect</Select.Option>
                        <Select.Option value={HostRouteType.STATIC_RESPONSE}>Static response</Select.Option>
                    </Select>
                </Form.Item>
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-route-source-path"
                    layout="vertical"
                    name={[name, "sourcePath"]}
                    validateStatus={validationResult.getStatus(`routes[${index}].sourcePath`)}
                    help={validationResult.getMessage(`routes[${index}].sourcePath`)}
                    label="Source path"
                    required
                >
                    <Input />
                </Form.Item>

                <If condition={type === HostRouteType.PROXY}>{this.renderProxyRoute(field, index)}</If>
                <If condition={type === HostRouteType.REDIRECT}>{this.renderRedirectRoute(field, index)}</If>
                <If condition={type === HostRouteType.STATIC_RESPONSE}>
                    {this.renderStaticResponseRoute(field, index)}
                </If>

                <Modal
                    title="Advanced settings"
                    open={index === customSettingsOpenModalIndex}
                    onClose={() => this.closeCustomSettingsModal()}
                    onCancel={() => this.closeCustomSettingsModal()}
                    footer={null}
                >
                    <p>
                        Any instruction placed here will be placed in the nginx configuration files as-is. Use this
                        field for any customized configuration parameters that you need in the host route.
                    </p>
                    <p>
                        Please note that the text below must be in the syntax expected by the nginx. Please refer to the
                        documentation at &nbsp;
                        <a
                            href="https://nginx.org/en/docs/http/ngx_http_core_module.html#location"
                            target="_blank"
                            rel="noreferrer"
                        >
                            this link
                        </a>
                        &nbsp; for more details. If you isn't sure about what to place here, it's probably the best to
                        leave it empty.
                    </p>

                    <Form.Item
                        {...FormLayout.ExpandedLabeledItem}
                        className="host-form-route-custom-settings"
                        name={[name, "customSettings"]}
                        validateStatus={validationResult.getStatus(`routes[${index}].customSettings`)}
                        help={validationResult.getMessage(`routes[${index}].customSettings`)}
                        required
                    >
                        <TextArea rows={10} />
                    </Form.Item>
                </Modal>

                <ArrowUpOutlined
                    onClick={() => this.moveRoute(operations, index, -1)}
                    style={index === 0 ? DISABLED_ACTION_ICON_STYLE : ENABLED_ACTION_ICON_STYLE}
                />
                <ArrowDownOutlined
                    onClick={() => this.moveRoute(operations, index, 1)}
                    style={index === routes.length - 1 ? DISABLED_ACTION_ICON_STYLE : ENABLED_ACTION_ICON_STYLE}
                />
                <SettingOutlined
                    onClick={() => this.openCustomSettingsModal(index)}
                    style={ENABLED_ACTION_ICON_STYLE}
                />

                <If condition={routes.length > 1}>
                    <CloseOutlined onClick={() => operations.remove(field.name)} style={ACTION_ICON_STYLE} />
                </If>
            </Flex>
        )
    }

    private renderRoutes(fields: FormListFieldData[], operations: FormListOperation) {
        const bindings = fields.map((field, index) => this.renderRoute(field, operations, index))

        const addAction = (
            <Form.Item>
                <Button
                    type="dashed"
                    onClick={() =>
                        operations.add({
                            ...DEFAULT_VALUES,
                            priority: fields.length,
                        })
                    }
                    icon={<PlusOutlined />}
                >
                    Add route
                </Button>
            </Form.Item>
        )

        return [...bindings, addAction]
    }

    render() {
        return <Form.List name="routes">{(fields, operations) => this.renderRoutes(fields, operations)}</Form.List>
    }
}
