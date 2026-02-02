import React from "react"
import { Flex, Form, InputNumber, Select, Space, Switch } from "antd"
import ValidationResult from "../../../core/validation/ValidationResult"
import { INTEGER_MAX } from "../SettingsConstants"
import { TimeUnit } from "../model/SettingsDto"
import { I18n } from "../../../core/i18n/I18n"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"

export interface NginxIgnitionSettingsTabProps {
    validationResult: ValidationResult
}

export default class NginxIgnitionSettingsTab extends React.Component<NginxIgnitionSettingsTabProps> {
    private renderExecutionIntervalFieldset(pathPrefix: string) {
        const { validationResult } = this.props
        return (
            <Form.Item label={<I18n id={MessageKey.FrontendSettingsTabsIgnitionExecutionInterval} />} required>
                <Space.Compact className="settings-form-input-wide">
                    <Form.Item
                        name={[pathPrefix, "intervalUnitCount"]}
                        validateStatus={validationResult.getStatus(`${pathPrefix}.intervalUnitCount`)}
                        help={validationResult.getMessage(`${pathPrefix}.intervalUnitCount`)}
                        noStyle
                    >
                        <InputNumber min={1} max={INTEGER_MAX} className="settings-form-input-wide" />
                    </Form.Item>
                    <Form.Item
                        name={[pathPrefix, "intervalUnit"]}
                        validateStatus={validationResult.getStatus(`${pathPrefix}.intervalUnit`)}
                        help={validationResult.getMessage(`${pathPrefix}.intervalUnit`)}
                        noStyle
                    >
                        <Select>
                            <Select.Option value={TimeUnit.DAYS}>
                                <I18n id={MessageKey.FrontendSettingsTabsIgnitionTimeUnitDays} />
                            </Select.Option>
                            <Select.Option value={TimeUnit.HOURS}>
                                <I18n id={MessageKey.FrontendSettingsTabsIgnitionTimeUnitHours} />
                            </Select.Option>
                            <Select.Option value={TimeUnit.MINUTES}>
                                <I18n id={MessageKey.FrontendSettingsTabsIgnitionTimeUnitMinutes} />
                            </Select.Option>
                        </Select>
                    </Form.Item>
                </Space.Compact>
            </Form.Item>
        )
    }

    render() {
        const { validationResult } = this.props
        return (
            <>
                <h2 className="settings-form-section-name">
                    <I18n id={MessageKey.FrontendSettingsTabsIgnitionSectionScheduledTasks} />
                </h2>
                <p className="settings-form-section-help-text">
                    <I18n id={MessageKey.FrontendSettingsTabsIgnitionSectionScheduledTasksHelp} />
                </p>

                <Flex className="settings-form-inner-flex-container">
                    <Flex className="settings-form-inner-flex-container-column settings-form-expanded-label-size">
                        <h3 className="settings-form-subsection-name">
                            <I18n id={MessageKey.FrontendSettingsTabsIgnitionLogRotation} />
                        </h3>
                        <Form.Item
                            name={["logRotation", "enabled"]}
                            validateStatus={validationResult.getStatus("logRotation.enabled")}
                            help={validationResult.getMessage("logRotation.enabled")}
                            label={<I18n id={MessageKey.FrontendSettingsTabsIgnitionAutoRotationEnabled} />}
                            required
                        >
                            <Switch />
                        </Form.Item>
                        <Form.Item
                            name={["logRotation", "maximumLines"]}
                            validateStatus={validationResult.getStatus("logRotation.maximumLines")}
                            help={validationResult.getMessage("logRotation.maximumLines")}
                            label={<I18n id={MessageKey.FrontendSettingsTabsIgnitionLinesToKeep} />}
                            required
                        >
                            <InputNumber min={0} max={99_999} className="settings-form-input-wide" />
                        </Form.Item>
                        {this.renderExecutionIntervalFieldset("logRotation")}
                    </Flex>
                    <Flex className="settings-form-inner-flex-container-column settings-form-expanded-label-size">
                        <h3 className="settings-form-subsection-name">
                            <I18n id={MessageKey.FrontendSettingsTabsIgnitionSslAutoRenew} />
                        </h3>
                        <Form.Item
                            name={["certificateAutoRenew", "enabled"]}
                            validateStatus={validationResult.getStatus("certificateAutoRenew.enabled")}
                            help={validationResult.getMessage("certificateAutoRenew.enabled")}
                            label={<I18n id={MessageKey.FrontendSettingsTabsIgnitionAutoRenewEnabled} />}
                            required
                        >
                            <Switch />
                        </Form.Item>
                        {this.renderExecutionIntervalFieldset("certificateAutoRenew")}
                    </Flex>
                </Flex>
            </>
        )
    }
}
