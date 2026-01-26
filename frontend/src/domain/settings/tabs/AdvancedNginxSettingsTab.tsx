import React from "react"
import SettingsFormValues from "../model/SettingsFormValues"
import ValidationResult from "../../../core/validation/ValidationResult"
import { Alert, Flex, Form, InputNumber, Space } from "antd"
import { INTEGER_MAX } from "../SettingsConstants"
import { Link } from "react-router-dom"
import FormLayout from "../../../core/components/form/FormLayout"
import TextArea from "antd/es/input/TextArea"
import { I18n } from "../../../core/i18n/I18n"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"

export interface AdvancedSettingsTabProps {
    formValues: SettingsFormValues
    validationResult: ValidationResult
}

export default class AdvancedNginxSettingsTab extends React.Component<AdvancedSettingsTabProps> {
    private renderCustomSettings() {
        const { validationResult } = this.props

        return (
            <>
                <h2 className="settings-form-section-name">
                    <I18n id={MessageKey.CommonCustomSettings} />
                </h2>
                <p className="settings-form-section-help-text">
                    <I18n id={MessageKey.FrontendSettingsTabsAdvancedCustomSettingsHelpPrefix} />{" "}
                    <Link
                        to="https://nginx.org/en/docs/http/ngx_http_core_module.html"
                        target="_blank"
                        rel="noreferrer"
                    >
                        <I18n id={MessageKey.CommonNginxDocLink} />
                    </Link>{" "}
                    <I18n id={MessageKey.FrontendSettingsTabsAdvancedCustomSettingsHelpSuffix} />
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

    private renderBufferSizeFieldset(property: string, label: React.ReactNode) {
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
                    <Space.Addon>
                        <I18n id={MessageKey.FrontendSettingsTabsAdvancedBufferSlotsOf} />
                    </Space.Addon>
                    <Form.Item
                        name={["nginx", "buffers", property, "sizeKb"]}
                        validateStatus={validationResult.getStatus(`nginx.buffers.${property}.sizeKb`)}
                        help={validationResult.getMessage(`nginx.buffers.${property}.sizeKb`)}
                        noStyle
                    >
                        <InputNumber min={1} max={INTEGER_MAX} />
                    </Form.Item>
                    <Space.Addon>
                        <I18n id={MessageKey.CommonUnitKbEach} />
                    </Space.Addon>
                </Space.Compact>
            </Form.Item>
        )
    }

    private renderBuffersSettings() {
        const { validationResult } = this.props
        return (
            <>
                <h2 className="settings-form-section-name">
                    <I18n id={MessageKey.FrontendSettingsTabsAdvancedSectionBuffers} />
                </h2>
                <p className="settings-form-section-help-text">
                    <I18n id={MessageKey.FrontendSettingsTabsAdvancedSectionBuffersHelp} />
                </p>
                <Flex className="settings-form-inner-flex-container">
                    <Flex className="settings-form-inner-flex-container-column settings-form-expanded-label-size">
                        <Form.Item
                            label={<I18n id={MessageKey.FrontendSettingsTabsAdvancedBufferClientBody} />}
                            required
                        >
                            <Space.Compact className="settings-form-input-wide">
                                <Form.Item
                                    name={["nginx", "buffers", "clientBodyKb"]}
                                    validateStatus={validationResult.getStatus("nginx.buffers.clientBodyKb")}
                                    help={validationResult.getMessage("nginx.buffers.clientBodyKb")}
                                    noStyle
                                >
                                    <InputNumber min={1} max={INTEGER_MAX} className="settings-form-input-wide" />
                                </Form.Item>
                                <Space.Addon>
                                    <I18n id={MessageKey.CommonUnitKb} />
                                </Space.Addon>
                            </Space.Compact>
                        </Form.Item>
                        <Form.Item
                            label={<I18n id={MessageKey.FrontendSettingsTabsAdvancedBufferClientHeaders} />}
                            required
                        >
                            <Space.Compact className="settings-form-input-wide">
                                <Form.Item
                                    name={["nginx", "buffers", "clientHeaderKb"]}
                                    validateStatus={validationResult.getStatus("nginx.buffers.clientHeaderKb")}
                                    help={validationResult.getMessage("nginx.buffers.clientHeaderKb")}
                                    noStyle
                                >
                                    <InputNumber min={1} max={INTEGER_MAX} className="settings-form-input-wide" />
                                </Form.Item>
                                <Space.Addon>
                                    <I18n id={MessageKey.CommonUnitKb} />
                                </Space.Addon>
                            </Space.Compact>
                        </Form.Item>
                    </Flex>
                    <Flex className="settings-form-inner-flex-container-column settings-form-expanded-label-size">
                        {this.renderBufferSizeFieldset(
                            "largeClientHeader",
                            <I18n id={MessageKey.FrontendSettingsTabsAdvancedBufferLargeClientHeader} />,
                        )}
                        {this.renderBufferSizeFieldset(
                            "output",
                            <I18n id={MessageKey.FrontendSettingsTabsAdvancedBufferOutput} />,
                        )}
                    </Flex>
                </Flex>
            </>
        )
    }

    render() {
        return (
            <>
                <Alert
                    message={<I18n id={MessageKey.CommonWarningProceedWithCaution} />}
                    description={<I18n id={MessageKey.CommonWarningProceedWithCautionDescription} />}
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
