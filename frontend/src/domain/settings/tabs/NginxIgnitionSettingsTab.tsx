import React from "react"
import { Flex, Form, InputNumber, Select, Space, Switch } from "antd"
import ValidationResult from "../../../core/validation/ValidationResult"
import { INTEGER_MAX } from "../SettingsConstants"
import { TimeUnit } from "../model/SettingsDto"

export interface NginxIgnitionSettingsTabProps {
    validationResult: ValidationResult
}

export default class NginxIgnitionSettingsTab extends React.Component<NginxIgnitionSettingsTabProps> {
    private renderExecutionIntervalFieldset(pathPrefix: string) {
        const { validationResult } = this.props
        return (
            <Form.Item label="Execution interval" required>
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
                            <Select.Option value={TimeUnit.DAYS}>days</Select.Option>
                            <Select.Option value={TimeUnit.HOURS}>hours</Select.Option>
                            <Select.Option value={TimeUnit.MINUTES}>minutes</Select.Option>
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
                <h2 className="settings-form-section-name" style={{ marginTop: 0 }}>
                    Scheduled tasks
                </h2>
                <p className="settings-form-section-help-text">Definition of the nginx ignition's housekeeping tasks</p>

                <Flex className="settings-form-inner-flex-container">
                    <Flex className="settings-form-inner-flex-container-column settings-form-expanded-label-size">
                        <h3 className="settings-form-subsection-name">Log rotation</h3>
                        <Form.Item
                            name={["logRotation", "enabled"]}
                            validateStatus={validationResult.getStatus("logRotation.enabled")}
                            help={validationResult.getMessage("logRotation.enabled")}
                            label="Auto rotation enabled"
                            required
                        >
                            <Switch />
                        </Form.Item>
                        <Form.Item
                            name={["logRotation", "maximumLines"]}
                            validateStatus={validationResult.getStatus("logRotation.maximumLines")}
                            help={validationResult.getMessage("logRotation.maximumLines")}
                            label="Lines to keep"
                            required
                        >
                            <InputNumber min={0} max={10000} className="settings-form-input-wide" />
                        </Form.Item>
                        {this.renderExecutionIntervalFieldset("logRotation")}
                    </Flex>
                    <Flex className="settings-form-inner-flex-container-column settings-form-expanded-label-size">
                        <h3 className="settings-form-subsection-name">SSL certificates auto renew</h3>
                        <Form.Item
                            name={["certificateAutoRenew", "enabled"]}
                            validateStatus={validationResult.getStatus("certificateAutoRenew.enabled")}
                            help={validationResult.getMessage("certificateAutoRenew.enabled")}
                            label="Auto renew enabled"
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
