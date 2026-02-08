import React from "react"
import { Flex, Form, Input, InputNumber, Select, Space, Switch } from "antd"
import HostBindings from "../../host/components/HostBindings"
import ValidationResult from "../../../core/validation/ValidationResult"
import SettingsFormValues from "../model/SettingsFormValues"
import { LogLevel } from "../model/SettingsDto"
import { INTEGER_MAX } from "../SettingsConstants"
import { I18n } from "../../../core/i18n/I18n"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { QuestionCircleFilled } from "@ant-design/icons"

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
                    label={<I18n id={MessageKey.FrontendSettingsTabsNginxErrorLogsEnabled} />}
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name={["nginx", "logs", levelField]}
                    validateStatus={validationResult.getStatus(`nginx.logs.${levelField}`)}
                    help={validationResult.getMessage(`nginx.logs.${levelField}`)}
                    label={<I18n id={MessageKey.FrontendSettingsTabsNginxErrorLogsLevel} />}
                    required
                >
                    <Select>
                        <Select.Option value={LogLevel.WARN}>
                            <I18n id={MessageKey.FrontendSettingsTabsNginxLogLevelWarn} />
                        </Select.Option>
                        <Select.Option value={LogLevel.ALERT}>
                            <I18n id={MessageKey.FrontendSettingsTabsNginxLogLevelAlert} />
                        </Select.Option>
                        <Select.Option value={LogLevel.ERROR}>
                            <I18n id={MessageKey.FrontendSettingsTabsNginxLogLevelError} />
                        </Select.Option>
                        <Select.Option value={LogLevel.CRIT}>
                            <I18n id={MessageKey.FrontendSettingsTabsNginxLogLevelCrit} />
                        </Select.Option>
                        <Select.Option value={LogLevel.EMERG}>
                            <I18n id={MessageKey.FrontendSettingsTabsNginxLogLevelEmerg} />
                        </Select.Option>
                    </Select>
                </Form.Item>
            </>
        )
    }

    private renderStatsColumn() {
        const { validationResult } = this.props
        return (
            <Flex className="settings-form-inner-flex-container-column settings-form-expanded-label-size">
                <h3 className="settings-form-subsection-name">
                    <I18n id={MessageKey.CommonTrafficStats} />
                </h3>
                <Form.Item
                    name={["nginx", "stats", "enabled"]}
                    validateStatus={validationResult.getStatus("nginx.stats.enabled")}
                    help={validationResult.getMessage("nginx.stats.enabled")}
                    label={<I18n id={MessageKey.FrontendSettingsTabsNginxStatsEnabled} />}
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name={["nginx", "stats", "persistent"]}
                    validateStatus={validationResult.getStatus("nginx.stats.persistent")}
                    help={validationResult.getMessage("nginx.stats.persistent")}
                    label={<I18n id={MessageKey.FrontendSettingsTabsNginxStatsPersistent} />}
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name={["nginx", "stats", "allHosts"]}
                    validateStatus={validationResult.getStatus("nginx.stats.allHosts")}
                    help={validationResult.getMessage("nginx.stats.allHosts")}
                    label={<I18n id={MessageKey.FrontendSettingsTabsNginxStatsAllHosts} />}
                    tooltip={{
                        title: <I18n id={MessageKey.FrontendSettingsTabsNginxStatsAllHostsHelp} />,
                        icon: <QuestionCircleFilled />,
                    }}
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item label={<I18n id={MessageKey.FrontendSettingsTabsNginxStatsMaximumSize} />} required>
                    <Space.Compact className="settings-form-input-wide">
                        <Form.Item
                            name={["nginx", "stats", "maximumSizeMb"]}
                            validateStatus={validationResult.getStatus("nginx.stats.maximumSizeMb")}
                            help={validationResult.getMessage("nginx.stats.maximumSizeMb")}
                            noStyle
                        >
                            <InputNumber min={1} max={512} className="settings-form-input-wide" />
                        </Form.Item>
                        <Space.Addon>
                            <I18n id={MessageKey.CommonUnitMb} />
                        </Space.Addon>
                    </Space.Compact>
                </Form.Item>
                <Form.Item
                    name={["nginx", "stats", "databaseLocation"]}
                    validateStatus={validationResult.getStatus("nginx.stats.databaseLocation")}
                    help={validationResult.getMessage("nginx.stats.databaseLocation")}
                    label={<I18n id={MessageKey.FrontendSettingsTabsNginxStatsDatabaseLocation} />}
                >
                    <Input maxLength={256} />
                </Form.Item>
            </Flex>
        )
    }

    private renderLogsFormPortion() {
        const { validationResult } = this.props
        return (
            <>
                <h2 className="settings-form-section-name">
                    <I18n id={MessageKey.CommonLogs} />
                </h2>
                <p className="settings-form-section-help-text">
                    <I18n id={MessageKey.FrontendSettingsTabsNginxLogsHelp} />
                </p>

                <Flex className="settings-form-inner-flex-container">
                    <Flex className="settings-form-inner-flex-container-column settings-form-expanded-label-size">
                        <h3 className="settings-form-subsection-name">
                            <I18n id={MessageKey.FrontendSettingsTabsNginxLogsVirtualHosts} />
                        </h3>
                        <Form.Item
                            name={["nginx", "logs", "accessLogsEnabled"]}
                            validateStatus={validationResult.getStatus("nginx.logs.accessLogsEnabled")}
                            help={validationResult.getMessage("nginx.logs.accessLogsEnabled")}
                            label={<I18n id={MessageKey.FrontendSettingsTabsNginxAccessLogsEnabled} />}
                            required
                        >
                            <Switch />
                        </Form.Item>
                        {this.renderErrorLogFieldset("errorLogsEnabled", "errorLogsLevel")}
                        <h3 className="settings-form-subsection-name">
                            <I18n id={MessageKey.FrontendSettingsTabsNginxLogsServer} />
                        </h3>
                        {this.renderErrorLogFieldset("serverLogsEnabled", "serverLogsLevel")}
                    </Flex>

                    {this.renderStatsColumn()}
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
                    label={<I18n id={MessageKey.FrontendSettingsTabsNginxConnectionsPerWorker} />}
                    required
                >
                    <InputNumber min={32} max={4096} className="settings-form-input-wide" />
                </Form.Item>
                <Form.Item label={<I18n id={MessageKey.FrontendSettingsTabsNginxMaximumBodySize} />} required>
                    <Space.Compact className="settings-form-input-wide">
                        <Form.Item
                            name={["nginx", "maximumBodySizeMb"]}
                            validateStatus={validationResult.getStatus("nginx.maximumBodySizeMb")}
                            help={validationResult.getMessage("nginx.maximumBodySizeMb")}
                            noStyle
                        >
                            <InputNumber min={1} max={INTEGER_MAX} className="settings-form-input-wide" />
                        </Form.Item>
                        <Space.Addon>
                            <I18n id={MessageKey.CommonUnitMb} />
                        </Space.Addon>
                    </Space.Compact>
                </Form.Item>
                <Form.Item label={<I18n id={MessageKey.FrontendSettingsTabsNginxConnectTimeout} />} required>
                    <Space.Compact className="settings-form-input-wide">
                        <Form.Item
                            name={["nginx", "timeouts", "connect"]}
                            validateStatus={validationResult.getStatus("nginx.timeouts.connect")}
                            help={validationResult.getMessage("nginx.timeouts.connect")}
                            noStyle
                        >
                            <InputNumber min={1} max={INTEGER_MAX} className="settings-form-input-wide" />
                        </Form.Item>
                        <Space.Addon>
                            <I18n id={MessageKey.CommonUnitSeconds} />
                        </Space.Addon>
                    </Space.Compact>
                </Form.Item>
                <Form.Item label={<I18n id={MessageKey.FrontendSettingsTabsNginxReadTimeout} />} required>
                    <Space.Compact className="settings-form-input-wide">
                        <Form.Item
                            name={["nginx", "timeouts", "read"]}
                            validateStatus={validationResult.getStatus("nginx.timeouts.read")}
                            help={validationResult.getMessage("nginx.timeouts.read")}
                            noStyle
                        >
                            <InputNumber min={1} max={INTEGER_MAX} className="settings-form-input-wide" />
                        </Form.Item>
                        <Space.Addon>
                            <I18n id={MessageKey.CommonUnitSeconds} />
                        </Space.Addon>
                    </Space.Compact>
                </Form.Item>
                <Form.Item label={<I18n id={MessageKey.FrontendSettingsTabsNginxSendTimeout} />} required>
                    <Space.Compact className="settings-form-input-wide">
                        <Form.Item
                            name={["nginx", "timeouts", "send"]}
                            validateStatus={validationResult.getStatus("nginx.timeouts.send")}
                            help={validationResult.getMessage("nginx.timeouts.send")}
                            noStyle
                        >
                            <InputNumber min={1} max={INTEGER_MAX} className="settings-form-input-wide" />
                        </Form.Item>
                        <Space.Addon>
                            <I18n id={MessageKey.CommonUnitSeconds} />
                        </Space.Addon>
                    </Space.Compact>
                </Form.Item>
                <Form.Item label={<I18n id={MessageKey.FrontendSettingsTabsNginxKeepaliveTimeout} />} required>
                    <Space.Compact className="settings-form-input-wide">
                        <Form.Item
                            name={["nginx", "timeouts", "keepalive"]}
                            validateStatus={validationResult.getStatus("nginx.timeouts.keepalive")}
                            help={validationResult.getMessage("nginx.timeouts.keepalive")}
                            noStyle
                        >
                            <InputNumber min={1} max={INTEGER_MAX} className="settings-form-input-wide" />
                        </Form.Item>
                        <Space.Addon>
                            <I18n id={MessageKey.CommonUnitSeconds} />
                        </Space.Addon>
                    </Space.Compact>
                </Form.Item>
                <Form.Item label={<I18n id={MessageKey.FrontendSettingsTabsNginxClientBodyTimeout} />} required>
                    <Space.Compact className="settings-form-input-wide">
                        <Form.Item
                            name={["nginx", "timeouts", "clientBody"]}
                            validateStatus={validationResult.getStatus("nginx.timeouts.clientBody")}
                            help={validationResult.getMessage("nginx.timeouts.clientBody")}
                            noStyle
                        >
                            <InputNumber min={1} max={INTEGER_MAX} className="settings-form-input-wide" />
                        </Form.Item>
                        <Space.Addon>
                            <I18n id={MessageKey.CommonUnitSeconds} />
                        </Space.Addon>
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
                    label={<I18n id={MessageKey.FrontendSettingsTabsNginxServerTokens} />}
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name={["nginx", "gzipEnabled"]}
                    validateStatus={validationResult.getStatus("nginx.gzipEnabled")}
                    help={validationResult.getMessage("nginx.gzipEnabled")}
                    label={<I18n id={MessageKey.FrontendSettingsTabsNginxGzipEnabled} />}
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name={["nginx", "sendfileEnabled"]}
                    validateStatus={validationResult.getStatus("nginx.sendfileEnabled")}
                    help={validationResult.getMessage("nginx.sendfileEnabled")}
                    label={<I18n id={MessageKey.FrontendSettingsTabsNginxSendfileEnabled} />}
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name={["nginx", "tcpNoDelayEnabled"]}
                    validateStatus={validationResult.getStatus("nginx.tcpNoDelayEnabled")}
                    help={validationResult.getMessage("nginx.tcpNoDelayEnabled")}
                    label={<I18n id={MessageKey.FrontendSettingsTabsNginxTcpNodelayEnabled} />}
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name={["nginx", "defaultContentType"]}
                    validateStatus={validationResult.getStatus("nginx.defaultContentType")}
                    help={validationResult.getMessage("nginx.defaultContentType")}
                    label={<I18n id={MessageKey.FrontendSettingsTabsNginxDefaultContentType} />}
                    required
                >
                    <Input maxLength={128} minLength={1} />
                </Form.Item>
                <Form.Item
                    name={["nginx", "runtimeUser"]}
                    validateStatus={validationResult.getStatus("nginx.runtimeUser")}
                    help={validationResult.getMessage("nginx.runtimeUser")}
                    label={<I18n id={MessageKey.FrontendSettingsTabsNginxRuntimeUser} />}
                    required
                >
                    <Input maxLength={32} minLength={1} />
                </Form.Item>
                <Form.Item
                    name={["nginx", "workerProcesses"]}
                    validateStatus={validationResult.getStatus("nginx.workerProcesses")}
                    help={validationResult.getMessage("nginx.workerProcesses")}
                    label={<I18n id={MessageKey.FrontendSettingsTabsNginxWorkerProcesses} />}
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
                <h2 className="settings-form-section-name">
                    <I18n id={MessageKey.CommonGeneral} />
                </h2>
                <p className="settings-form-section-help-text">
                    <I18n id={MessageKey.FrontendSettingsTabsNginxSectionGeneralHelp} />
                </p>
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
                <h2 className="settings-form-section-name">
                    <I18n id={MessageKey.FrontendSettingsTabsNginxSectionBindings} />
                </h2>
                <p className="settings-form-section-help-text">
                    <I18n id={MessageKey.FrontendSettingsTabsNginxSectionBindingsHelp} />
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
