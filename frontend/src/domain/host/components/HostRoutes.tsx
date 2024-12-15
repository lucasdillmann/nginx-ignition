import React from "react"
import ValidationResult from "../../../core/validation/ValidationResult"
import { Button, Flex, Form, FormListFieldData, FormListOperation, Input, InputNumber, Select } from "antd"
import If from "../../../core/components/flowcontrol/If"
import { CloseOutlined, PlusOutlined, ArrowUpOutlined, ArrowDownOutlined, SettingOutlined } from "@ant-design/icons"
import { HostFormRoute } from "../model/HostFormValues"
import { HostRouteType } from "../model/HostRequest"
import FormLayout from "../../../core/components/form/FormLayout"
import "./HostRoutes.css"
import TextArea from "antd/es/input/TextArea"
import { IntegrationResponse } from "../../integration/model/IntegrationResponse"
import PaginatedSelect from "../../../core/components/select/PaginatedSelect"
import { IntegrationOptionResponse } from "../../integration/model/IntegrationOptionResponse"
import PageResponse, { emptyPageResponse } from "../../../core/pagination/PageResponse"
import IntegrationService from "../../integration/IntegrationService"
import HostFormValuesDefaults from "../model/HostFormValuesDefaults"
import HostRouteSettingsModal from "./HostRouteSettingsModal"

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

export interface HostRoutesProps {
    routes: HostFormRoute[]
    validationResult: ValidationResult
    integrations: IntegrationResponse[]
    onRouteRemove: (index: number) => void
}

interface HostRoutesState {
    routeSettingsOpenModalIndex?: number
}

export default class HostRoutes extends React.Component<HostRoutesProps, HostRoutesState> {
    private readonly integrationService: IntegrationService
    private readonly optionsRef: React.RefObject<PaginatedSelect<IntegrationOptionResponse>>

    constructor(props: HostRoutesProps) {
        super(props)
        this.integrationService = new IntegrationService()
        this.optionsRef = React.createRef()
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

    private async fetchIntegrationOptions(
        pageSize: number,
        pageNumber: number,
        searchTerms?: string,
        integrationId?: string,
    ): Promise<PageResponse<IntegrationOptionResponse>> {
        if (integrationId === undefined) return emptyPageResponse<IntegrationOptionResponse>()

        return this.integrationService.getOptions(integrationId!!, pageSize, pageNumber, searchTerms)
    }

    private handleIntegrationChange() {
        this.optionsRef.current?.reset()
    }

    private renderIntegrationRoute(field: FormListFieldData, index: number): React.ReactNode {
        const { validationResult, integrations, routes } = this.props
        const { name } = field
        const integrationOptions = integrations.map(({ id, name }) => ({
            label: name,
            value: id,
        }))

        const currentIntegrationId = routes[index].integration?.integrationId
        return (
            <>
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-route-integration-app"
                    layout="vertical"
                    name={[name, "integration", "integrationId"]}
                    validateStatus={validationResult.getStatus(`routes[${index}].integration.integrationId`)}
                    help={validationResult.getMessage(`routes[${index}].integration.integrationId`)}
                    label="Integration"
                    required
                >
                    <Select options={integrationOptions} onChange={() => this.handleIntegrationChange()} />
                </Form.Item>
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-route-integration-option"
                    layout="vertical"
                    name={[name, "integration", "option"]}
                    validateStatus={validationResult.getStatus(`routes[${index}].integration.optionId`)}
                    help={validationResult.getMessage(`routes[${index}].integration.optionId`)}
                    label="Option / App"
                    required
                >
                    <PaginatedSelect<IntegrationOptionResponse>
                        ref={this.optionsRef}
                        disabled={currentIntegrationId === undefined}
                        itemKey={item => item.id}
                        itemDescription={item => item.name}
                        pageProvider={(pageSize, pageNumber, searchTerms) =>
                            this.fetchIntegrationOptions(pageSize, pageNumber, searchTerms, currentIntegrationId)
                        }
                    />
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

    private openRouteSettingsModal(index: number) {
        this.setState({ routeSettingsOpenModalIndex: index })
    }

    private closeRouteSettingsModal() {
        this.setState({ routeSettingsOpenModalIndex: undefined })
    }

    private removeRoute(index: number) {
        const { onRouteRemove } = this.props
        onRouteRemove(index)
    }

    private renderRoute(field: FormListFieldData, operations: FormListOperation, index: number) {
        const { validationResult, routes } = this.props
        const { routeSettingsOpenModalIndex } = this.state
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
                        <Select.Option value={HostRouteType.INTEGRATION}>Integration</Select.Option>
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

                <If condition={type === HostRouteType.INTEGRATION}>{this.renderIntegrationRoute(field, index)}</If>
                <If condition={type === HostRouteType.PROXY}>{this.renderProxyRoute(field, index)}</If>
                <If condition={type === HostRouteType.REDIRECT}>{this.renderRedirectRoute(field, index)}</If>
                <If condition={type === HostRouteType.STATIC_RESPONSE}>
                    {this.renderStaticResponseRoute(field, index)}
                </If>

                <HostRouteSettingsModal
                    open={index === routeSettingsOpenModalIndex}
                    onClose={() => this.closeRouteSettingsModal()}
                    onCancel={() => this.closeRouteSettingsModal()}
                    index={index}
                    fieldPath={name}
                    validationResult={validationResult}
                />

                <ArrowUpOutlined
                    onClick={() => this.moveRoute(operations, index, -1)}
                    style={index === 0 ? DISABLED_ACTION_ICON_STYLE : ENABLED_ACTION_ICON_STYLE}
                />
                <ArrowDownOutlined
                    onClick={() => this.moveRoute(operations, index, 1)}
                    style={index === routes.length - 1 ? DISABLED_ACTION_ICON_STYLE : ENABLED_ACTION_ICON_STYLE}
                />
                <SettingOutlined onClick={() => this.openRouteSettingsModal(index)} style={ENABLED_ACTION_ICON_STYLE} />

                <If condition={routes.length > 1}>
                    <CloseOutlined onClick={() => this.removeRoute(index)} style={ACTION_ICON_STYLE} />
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
                            ...HostFormValuesDefaults.routes[0],
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
