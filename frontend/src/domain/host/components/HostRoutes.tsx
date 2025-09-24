import React from "react"
import ValidationResult from "../../../core/validation/ValidationResult"
import { Button, Flex, Form, FormListFieldData, FormListOperation, Input, InputNumber, Select, Switch } from "antd"
import If from "../../../core/components/flowcontrol/If"
import {
    ArrowDownOutlined,
    ArrowUpOutlined,
    DeleteOutlined,
    PlusOutlined,
    QuestionCircleFilled,
    SettingOutlined,
} from "@ant-design/icons"
import { HostFormRoute } from "../model/HostFormValues"
import { HostRouteSourceCodeLanguage, HostRouteType } from "../model/HostRequest"
import FormLayout from "../../../core/components/form/FormLayout"
import "./HostRoutes.css"
import TextArea from "antd/es/input/TextArea"
import { IntegrationResponse } from "../../integration/model/IntegrationResponse"
import PaginatedSelect from "../../../core/components/select/PaginatedSelect"
import { IntegrationOptionResponse } from "../../integration/model/IntegrationOptionResponse"
import PageResponse, { emptyPageResponse } from "../../../core/pagination/PageResponse"
import IntegrationService from "../../integration/IntegrationService"
import HostRouteSettingsModal from "./HostRouteSettingsModal"
import { Link } from "react-router-dom"
import CodeEditorModal from "../../../core/components/codeeditor/CodeEditorModal"
import { CodeEditorLanguage } from "../../../core/components/codeeditor/CodeEditor"
import { hostFormValuesDefaults } from "../model/HostFormValuesDefaults"

const ACTION_ICON_STYLE = {
    marginLeft: 15,
    alignItems: "start",
    marginTop: 37,
}

const DISABLED_ACTION_ICON_STYLE = {
    ...ACTION_ICON_STYLE,
    color: "var(--nginxIgnition-colorTextSecondary)",
}

const ENABLED_ACTION_ICON_STYLE = {
    ...ACTION_ICON_STYLE,
    color: "var(--nginxIgnition-colorText)",
}

export interface HostRoutesProps {
    routes: HostFormRoute[]
    validationResult: ValidationResult
    integrations: IntegrationResponse[]
    onRouteRemove: (index: number) => void
    onChange: () => void
}

interface HostRoutesState {
    routeSettingsOpenModalIndex?: number
    routeCodeEditorOpenModalIndex?: number
}

export default class HostRoutes extends React.Component<HostRoutesProps, HostRoutesState> {
    private readonly integrationService: IntegrationService
    private readonly optionsRef: React.RefObject<PaginatedSelect<IntegrationOptionResponse> | null>

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

    private renderStaticFilesRoute(field: FormListFieldData, index: number): React.ReactNode {
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
                label="Target directory path"
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
                    name={[name, "redirectCode"]}
                    validateStatus={validationResult.getStatus(`routes[${index}].redirectCode`)}
                    help={validationResult.getMessage(`routes[${index}].redirectCode`)}
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

        return this.integrationService.getOptions(integrationId!!, pageSize, pageNumber, searchTerms, true)
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
                        itemDescription={item => `${item.name} (${item.port} HTTP)`}
                        pageProvider={(pageSize, pageNumber, searchTerms) =>
                            this.fetchIntegrationOptions(pageSize, pageNumber, searchTerms, currentIntegrationId)
                        }
                    />
                </Form.Item>
            </>
        )
    }

    private handleRoutePayloadChange(index: number, payload: string) {
        const { routes } = this.props

        const route = routes[index]
        if (!route?.response) {
            route.response = {
                statusCode: 200,
            }
        }

        route.response.payload = payload
    }

    private handleRouteCodeChange(index: number, code: string) {
        const { routes } = this.props

        const route = routes[index]
        if (!route?.sourceCode) {
            return
        }

        route.sourceCode.code = code
    }

    private renderStaticResponseRoute(field: FormListFieldData, index: number): React.ReactNode {
        const { validationResult, routes } = this.props
        const { routeCodeEditorOpenModalIndex } = this.state
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
                    <TextArea rows={3} onClick={() => this.openRouteCodeEditorModal(index)} />
                </Form.Item>

                <CodeEditorModal
                    open={index === routeCodeEditorOpenModalIndex}
                    onClose={() => this.closeRouteCodeEditorModal()}
                    value={routes[index]?.response?.payload ?? ""}
                    onChange={value => this.handleRoutePayloadChange(index, value)}
                    language={[
                        CodeEditorLanguage.JSON,
                        CodeEditorLanguage.HTML,
                        CodeEditorLanguage.CSS,
                        CodeEditorLanguage.XML,
                        CodeEditorLanguage.YAML,
                        CodeEditorLanguage.JAVASCRIPT,
                        CodeEditorLanguage.PLAIN_TEXT,
                    ]}
                />
            </>
        )
    }

    private renderSourceCodeHelpText(index: number): React.ReactNode {
        const { routes } = this.props
        const { sourceCode } = routes[index]
        const url =
            sourceCode?.language === HostRouteSourceCodeLanguage.LUA
                ? "https://github.com/openresty/lua-nginx-module/blob/master/README.markdown#synopsis"
                : "https://nginx.org/en/docs/njs/"

        return (
            <>
                Check{" "}
                <Link to={url} target="_blank">
                    this link
                </Link>{" "}
                for instructions
            </>
        )
    }

    private renderSourceCodeRoute(field: FormListFieldData, index: number): React.ReactNode {
        const { validationResult, routes } = this.props
        const { routeCodeEditorOpenModalIndex } = this.state
        const { name } = field
        const { sourceCode } = routes[index]
        const helpText = this.renderSourceCodeHelpText(index)

        return (
            <>
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-route-source-code-language"
                    layout="vertical"
                    name={[name, "sourceCode", "language"]}
                    validateStatus={validationResult.getStatus(`routes[${index}].sourceCode.language`)}
                    help={validationResult.getMessage(`routes[${index}].sourceCode.language`)}
                    label="Language"
                    required
                >
                    <Select>
                        <Select.Option value={HostRouteSourceCodeLanguage.JAVASCRIPT}>JavaScript</Select.Option>
                        <Select.Option value={HostRouteSourceCodeLanguage.LUA}>Lua</Select.Option>
                    </Select>
                </Form.Item>
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-route-source-code-code"
                    layout="vertical"
                    name={[name, "sourceCode", "code"]}
                    validateStatus={validationResult.getStatus(`routes[${index}].sourceCode.code`)}
                    help={validationResult.getMessage(`routes[${index}].sourceCode.code`) ?? helpText}
                    label="Source code"
                    required
                >
                    <TextArea rows={3} onClick={() => this.openRouteCodeEditorModal(index)} />
                </Form.Item>
                <If condition={sourceCode?.language === HostRouteSourceCodeLanguage.JAVASCRIPT}>
                    <Form.Item
                        {...FormLayout.ExpandedLabeledItem}
                        className="host-form-route-source-code-main-function"
                        layout="vertical"
                        name={[name, "sourceCode", "mainFunction"]}
                        validateStatus={validationResult.getStatus(`routes[${index}].sourceCode.mainFunction`)}
                        help={validationResult.getMessage(`routes[${index}].sourceCode.mainFunction`)}
                        label="Main function name"
                        required
                    >
                        <Input />
                    </Form.Item>
                </If>

                <CodeEditorModal
                    open={index === routeCodeEditorOpenModalIndex}
                    onClose={() => this.closeRouteCodeEditorModal()}
                    value={routes[index]?.sourceCode?.code ?? ""}
                    onChange={value => this.handleRouteCodeChange(index, value)}
                    language={
                        sourceCode?.language === HostRouteSourceCodeLanguage.JAVASCRIPT
                            ? CodeEditorLanguage.JAVASCRIPT
                            : CodeEditorLanguage.LUA
                    }
                />
            </>
        )
    }

    private moveRoute(operations: FormListOperation, index: number, offset: number) {
        const { routes } = this.props

        if (index === 0 && offset < 0) return
        if (index === routes.length - 1 && offset > 0) return

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

    private openRouteCodeEditorModal(index: number) {
        this.setState({ routeCodeEditorOpenModalIndex: index })
    }

    private closeRouteCodeEditorModal() {
        const { onChange } = this.props

        this.setState({ routeCodeEditorOpenModalIndex: undefined })
        onChange()
    }

    private removeRoute(index: number) {
        const { onRouteRemove } = this.props
        onRouteRemove(index)
    }

    private buildRouteTypeTooltipContents() {
        return (
            <>
                <p>The route type defines how the requests should be handled</p>
                <p>
                    <b>Integration:</b> Proxies requests to an app running in a TrueNAS or Docker container
                </p>
                <p>
                    <b>Proxy:</b> Proxies requests to an app using a target URL
                </p>
                <p>
                    <b>Redirect:</b> Redirects requests to an third-party URL
                </p>
                <p>
                    <b>Static response:</b> Returns a static response with a predefined status code, headers and body
                </p>
                <p>
                    <b>Source code:</b> Executes a JavaScript or Lua code to handle requests
                </p>
                <p>
                    <b>Directory:</b> Serves static files from a directory with listing enabled
                </p>
            </>
        )
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
                    tooltip={{
                        title: this.buildRouteTypeTooltipContents(),
                        icon: <QuestionCircleFilled />,
                    }}
                    required
                >
                    <Select>
                        <Select.Option value={HostRouteType.INTEGRATION}>Integration</Select.Option>
                        <Select.Option value={HostRouteType.PROXY}>Proxy</Select.Option>
                        <Select.Option value={HostRouteType.REDIRECT}>Redirect</Select.Option>
                        <Select.Option value={HostRouteType.STATIC_RESPONSE}>Static response</Select.Option>
                        <Select.Option value={HostRouteType.STATIC_FILES}>Static files</Select.Option>
                        <Select.Option value={HostRouteType.EXECUTE_CODE}>Execute code</Select.Option>
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
                <If condition={type === HostRouteType.EXECUTE_CODE}>{this.renderSourceCodeRoute(field, index)}</If>
                <If condition={type === HostRouteType.STATIC_FILES}>{this.renderStaticFilesRoute(field, index)}</If>

                <HostRouteSettingsModal
                    open={index === routeSettingsOpenModalIndex}
                    route={routes[index]}
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

                <Form.Item style={{ ...ENABLED_ACTION_ICON_STYLE, marginTop: 28 }} name={[name, "enabled"]} required>
                    <Switch size="small" />
                </Form.Item>

                <If condition={routes.length > 1}>
                    <DeleteOutlined onClick={() => this.removeRoute(index)} style={ACTION_ICON_STYLE} />
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
                            ...hostFormValuesDefaults().routes[0],
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
