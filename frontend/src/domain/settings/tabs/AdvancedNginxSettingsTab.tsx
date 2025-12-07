import React from "react"
import SettingsFormValues from "../model/SettingsFormValues"
import ValidationResult from "../../../core/validation/ValidationResult"
import { Alert, Flex, Form, InputNumber, Space } from "antd"
import { INTEGER_MAX } from "../SettingsConstants"
import { Link } from "react-router-dom"
import FormLayout from "../../../core/components/form/FormLayout"
import TextArea from "antd/es/input/TextArea"

export interface AdvancedSettingsTabProps {
    formValues: SettingsFormValues
    validationResult: ValidationResult
}

export default class AdvancedNginxSettingsTab extends React.Component<AdvancedSettingsTabProps> {
    private renderCustomSettings() {
        const { validationResult } = this.props

        return (
            <>
                <h2 className="settings-form-section-name">Custom settings</h2>
                <p className="settings-form-section-help-text">
                    Anything placed here will be added as-is to the http directive of the main nginx configuration file.
                    Custom configuration must be in the syntax expected by the nginx (refer to the documentation
                    at&nbsp;
                    <Link
                        to="https://nginx.org/en/docs/http/ngx_http_core_module.html"
                        target="_blank"
                        rel="noreferrer"
                    >
                        this link
                    </Link>
                    &nbsp;for more details). If you isn't sure about what to place here, it's probably the best to leave
                    it empty.
                </p>

                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-route-custom-settings"
                    name={["nginx", "custom"]}
                    validateStatus={validationResult.getStatus("nginx.custom")}
                    help={validationResult.getMessage("nginx.custom")}
                    required
                >
                    <TextArea rows={10} />
                </Form.Item>
            </>
        )
    }

    private renderBufferSizeFieldset(property: string, label: string) {
        const { validationResult } = this.props
        return (
            <Form.Item label={label} required>
                <Space.Compact>
                    <Form.Item
                        name={["nginx", "buffers", property, "amount"]}
                        validateStatus={validationResult.getStatus(`nginx.buffers.${property}.amount`)}
                        help={validationResult.getMessage(`nginx.buffers.${property}.amount`)}
                        noStyle
                    >
                        <InputNumber min={1} max={INTEGER_MAX} />
                    </Form.Item>
                    <Space.Addon>slots of</Space.Addon>
                    <Form.Item
                        name={["nginx", "buffers", property, "sizeKb"]}
                        validateStatus={validationResult.getStatus(`nginx.buffers.${property}.sizeKb`)}
                        help={validationResult.getMessage(`nginx.buffers.${property}.sizeKb`)}
                        noStyle
                    >
                        <InputNumber min={1} max={INTEGER_MAX} />
                    </Form.Item>
                    <Space.Addon>KB each</Space.Addon>
                </Space.Compact>
            </Form.Item>
        )
    }

    private renderBuffersSettings() {
        const { validationResult } = this.props
        return (
            <>
                <h2 className="settings-form-section-name" style={{ marginTop: 0 }}>
                    Buffers
                </h2>
                <p className="settings-form-section-help-text">
                    Configuration of the nginx's buffering of the input/output traffic
                </p>
                <Flex className="settings-form-inner-flex-container">
                    <Flex className="settings-form-inner-flex-container-column settings-form-expanded-label-size">
                        <Form.Item label="Client body" required>
                            <Space.Compact className="settings-form-input-wide">
                                <Form.Item
                                    name={["nginx", "buffers", "clientBodyKb"]}
                                    validateStatus={validationResult.getStatus("nginx.buffers.clientBodyKb")}
                                    help={validationResult.getMessage("nginx.buffers.clientBodyKb")}
                                    noStyle
                                >
                                    <InputNumber min={1} max={INTEGER_MAX} className="settings-form-input-wide" />
                                </Form.Item>
                                <Space.Addon>KB</Space.Addon>
                            </Space.Compact>
                        </Form.Item>
                        <Form.Item label="Client headers" required>
                            <Space.Compact className="settings-form-input-wide">
                                <Form.Item
                                    name={["nginx", "buffers", "clientHeaderKb"]}
                                    validateStatus={validationResult.getStatus("nginx.buffers.clientHeaderKb")}
                                    help={validationResult.getMessage("nginx.buffers.clientHeaderKb")}
                                    noStyle
                                >
                                    <InputNumber min={1} max={INTEGER_MAX} className="settings-form-input-wide" />
                                </Form.Item>
                                <Space.Addon>KB</Space.Addon>
                            </Space.Compact>
                        </Form.Item>
                    </Flex>
                    <Flex className="settings-form-inner-flex-container-column settings-form-expanded-label-size">
                        {this.renderBufferSizeFieldset("largeClientHeader", "Large client headers")}
                        {this.renderBufferSizeFieldset("output", "Output")}
                    </Flex>
                </Flex>
            </>
        )
    }

    render() {
        return (
            <>
                <Alert
                    message="Proceed with caution"
                    description={`
                        These settings can break your nginx server or make it misbehave. Using the default values will 
                        work just fine for almost all use cases. 
                    `}
                    type="warning"
                    style={{ marginBottom: 20 }}
                    showIcon
                />
                {this.renderBuffersSettings()}
                {this.renderCustomSettings()}
            </>
        )
    }
}
