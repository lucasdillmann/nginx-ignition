import React from "react"
import { Flex, Form, Input, InputNumber, Select, Space, Switch } from "antd"
import HostBindings from "../../host/components/HostBindings"
import ValidationResult from "../../../core/validation/ValidationResult"
import SettingsFormValues from "../model/SettingsFormValues"
import { LogLevel, RuntimeUser } from "../model/SettingsDto"
import { INTEGER_MAX } from "../SettingsConstants"

export interface NginxSettingsTabProps {
    formValues: SettingsFormValues
    validationResult: ValidationResult
}

export default class NginxSettingsTab extends React.Component<NginxSettingsTabProps> {
    private renderErrorLogFieldset(enabledField: string, levelField: string) {
        const { validationResult } = this.props
        return (
            <>
                <Form.Item
                    name={["nginx", "logs", enabledField]}
                    validateStatus={validationResult.getStatus(`nginx.logs.${enabledField}`)}
                    help={validationResult.getMessage(`nginx.logs.${enabledField}`)}
                    label="Error logs enabled"
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name={["nginx", "logs", levelField]}
                    validateStatus={validationResult.getStatus(`nginx.logs.${levelField}`)}
                    help={validationResult.getMessage(`nginx.logs.${levelField}`)}
                    label="Level"
                    required
                >
                    <Select>
                        <Select.Option value={LogLevel.WARN}>warn</Select.Option>
                        <Select.Option value={LogLevel.ALERT}>alert</Select.Option>
                        <Select.Option value={LogLevel.ERROR}>error</Select.Option>
                        <Select.Option value={LogLevel.CRIT}>crit</Select.Option>
                        <Select.Option value={LogLevel.EMERG}>emerg</Select.Option>
                    </Select>
                </Form.Item>
            </>
        )
    }

    private renderLogsFormPortion() {
        const { validationResult } = this.props
        return (
            <>
                <h2 className="settings-form-section-name">Logs</h2>
                <p className="settings-form-section-help-text">
                    Logging settings for the nginx server and virtual hosts
                </p>

                <Flex className="settings-form-inner-flex-container">
                    <Flex className="settings-form-inner-flex-container-column settings-form-expanded-label-size">
                        <h3 className="settings-form-subsection-name">Virtual hosts</h3>
                        <Form.Item
                            name={["nginx", "logs", "accessLogsEnabled"]}
                            validateStatus={validationResult.getStatus("nginx.logs.accessLogsEnabled")}
                            help={validationResult.getMessage("nginx.logs.accessLogsEnabled")}
                            label="Access logs enabled"
                            required
                        >
                            <Switch />
                        </Form.Item>
                        {this.renderErrorLogFieldset("errorLogsEnabled", "errorLogsLevel")}
                    </Flex>

                    <Flex className="settings-form-inner-flex-container-column settings-form-expanded-label-size">
                        <h3 className="settings-form-subsection-name">Server</h3>
                        {this.renderErrorLogFieldset("serverLogsEnabled", "serverLogsLevel")}
                    </Flex>
                </Flex>
            </>
        )
    }

    private renderGeneralSecondColumn() {
        const { validationResult } = this.props
        return (
            <Flex className="settings-form-inner-flex-container-column settings-form-expanded-label-size">
                <Form.Item
                    name={["nginx", "workerConnections"]}
                    validateStatus={validationResult.getStatus("nginx.workerConnections")}
                    help={validationResult.getMessage("nginx.workerConnections")}
                    label="Connections per worker"
                    required
                >
                    <InputNumber min={32} max={4096} className="settings-form-input-wide" />
                </Form.Item>
                <Form.Item label="Maximum body size" required>
                    <Space.Compact className="settings-form-input-wide">
                        <Form.Item
                            name={["nginx", "maximumBodySizeMb"]}
                            validateStatus={validationResult.getStatus("nginx.maximumBodySizeMb")}
                            help={validationResult.getMessage("nginx.maximumBodySizeMb")}
                            noStyle
                        >
                            <InputNumber min={1} max={INTEGER_MAX} className="settings-form-input-wide" />
                        </Form.Item>
                        <Space.Addon>MB</Space.Addon>
                    </Space.Compact>
                </Form.Item>
                <Form.Item label="Connect timeout" required>
                    <Space.Compact className="settings-form-input-wide">
                        <Form.Item
                            name={["nginx", "timeouts", "connect"]}
                            validateStatus={validationResult.getStatus("nginx.timeouts.connect")}
                            help={validationResult.getMessage("nginx.timeouts.connect")}
                            noStyle
                        >
                            <InputNumber min={1} max={INTEGER_MAX} className="settings-form-input-wide" />
                        </Form.Item>
                        <Space.Addon>seconds</Space.Addon>
                    </Space.Compact>
                </Form.Item>
                <Form.Item label="Read timeout" required>
                    <Space.Compact className="settings-form-input-wide">
                        <Form.Item
                            name={["nginx", "timeouts", "read"]}
                            validateStatus={validationResult.getStatus("nginx.timeouts.read")}
                            help={validationResult.getMessage("nginx.timeouts.read")}
                            noStyle
                        >
                            <InputNumber min={1} max={INTEGER_MAX} className="settings-form-input-wide" />
                        </Form.Item>
                        <Space.Addon>seconds</Space.Addon>
                    </Space.Compact>
                </Form.Item>
                <Form.Item label="Send timeout" required>
                    <Space.Compact className="settings-form-input-wide">
                        <Form.Item
                            name={["nginx", "timeouts", "send"]}
                            validateStatus={validationResult.getStatus("nginx.timeouts.send")}
                            help={validationResult.getMessage("nginx.timeouts.send")}
                            noStyle
                        >
                            <InputNumber min={1} max={INTEGER_MAX} className="settings-form-input-wide" />
                        </Form.Item>
                        <Space.Addon>seconds</Space.Addon>
                    </Space.Compact>
                </Form.Item>
                <Form.Item label="Keepalive timeout" required>
                    <Space.Compact className="settings-form-input-wide">
                        <Form.Item
                            name={["nginx", "timeouts", "keepalive"]}
                            validateStatus={validationResult.getStatus("nginx.timeouts.keepalive")}
                            help={validationResult.getMessage("nginx.timeouts.keepalive")}
                            noStyle
                        >
                            <InputNumber min={1} max={INTEGER_MAX} className="settings-form-input-wide" />
                        </Form.Item>
                        <Space.Addon>seconds</Space.Addon>
                    </Space.Compact>
                </Form.Item>
                <Form.Item label="Client body timeout" required>
                    <Space.Compact className="settings-form-input-wide">
                        <Form.Item
                            name={["nginx", "timeouts", "clientBody"]}
                            validateStatus={validationResult.getStatus("nginx.timeouts.clientBody")}
                            help={validationResult.getMessage("nginx.timeouts.clientBody")}
                            noStyle
                        >
                            <InputNumber min={1} max={INTEGER_MAX} className="settings-form-input-wide" />
                        </Form.Item>
                        <Space.Addon>seconds</Space.Addon>
                    </Space.Compact>
                </Form.Item>
            </Flex>
        )
    }

    private renderGeneralFirstColumn() {
        const { validationResult } = this.props
        return (
            <Flex className="settings-form-inner-flex-container-column settings-form-expanded-label-size">
                <Form.Item
                    name={["nginx", "serverTokensEnabled"]}
                    validateStatus={validationResult.getStatus("nginx.serverTokensEnabled")}
                    help={validationResult.getMessage("nginx.serverTokensEnabled")}
                    label="Server tokens"
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name={["nginx", "gzipEnabled"]}
                    validateStatus={validationResult.getStatus("nginx.gzipEnabled")}
                    help={validationResult.getMessage("nginx.gzipEnabled")}
                    label="GZIP enabled"
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name={["nginx", "sendfileEnabled"]}
                    validateStatus={validationResult.getStatus("nginx.sendfileEnabled")}
                    help={validationResult.getMessage("nginx.sendfileEnabled")}
                    label="Sendfile enabled"
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name={["nginx", "tcpNoDelayEnabled"]}
                    validateStatus={validationResult.getStatus("nginx.tcpNoDelayEnabled")}
                    help={validationResult.getMessage("nginx.tcpNoDelayEnabled")}
                    label="TCP nodelay enabled"
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name={["nginx", "defaultContentType"]}
                    validateStatus={validationResult.getStatus("nginx.defaultContentType")}
                    help={validationResult.getMessage("nginx.defaultContentType")}
                    label="Default content type"
                    required
                >
                    <Input maxLength={128} minLength={1} />
                </Form.Item>
                <Form.Item
                    name={["nginx", "runtimeUser"]}
                    validateStatus={validationResult.getStatus("nginx.runtimeUser")}
                    help={validationResult.getMessage("nginx.runtimeUser")}
                    label="Runtime user"
                    required
                >
                    <Select>
                        <Select.Option value={RuntimeUser.ROOT}>root</Select.Option>
                        <Select.Option value={RuntimeUser.NGINX}>nginx</Select.Option>
                    </Select>
                </Form.Item>
                <Form.Item
                    name={["nginx", "workerProcesses"]}
                    validateStatus={validationResult.getStatus("nginx.workerProcesses")}
                    help={validationResult.getMessage("nginx.workerProcesses")}
                    label="Worker processes"
                    required
                >
                    <InputNumber min={1} max={100} className="settings-form-input-wide" />
                </Form.Item>
            </Flex>
        )
    }

    private renderGeneralFormPortion() {
        return (
            <>
                <h2 className="settings-form-section-name">General</h2>
                <p className="settings-form-section-help-text">General configurations properties of the nginx server</p>
                <Flex className="settings-form-inner-flex-container">
                    {this.renderGeneralFirstColumn()}
                    {this.renderGeneralSecondColumn()}
                </Flex>
            </>
        )
    }

    private renderBindingsFormPortion() {
        const { validationResult, formValues } = this.props

        return (
            <>
                <h2 className="settings-form-section-name">Global bindings</h2>
                <p className="settings-form-section-help-text">
                    Relation of IPs and ports where the virtual hosts will listen for requests by default (can be
                    overwritten on every host if needed)
                </p>
                <HostBindings
                    pathPrefix="globalBindings"
                    bindings={formValues.globalBindings}
                    validationResult={validationResult}
                />
            </>
        )
    }

    render() {
        return (
            <>
                {this.renderGeneralFormPortion()}
                {this.renderLogsFormPortion()}
                {this.renderBindingsFormPortion()}
            </>
        )
    }
}
