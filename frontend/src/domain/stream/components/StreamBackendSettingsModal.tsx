import React from "react"
import { Form, InputNumber, Modal, Switch } from "antd"
import ValidationResult from "../../../core/validation/ValidationResult"
import { StreamBackend } from "../model/StreamRequest"
import { streamCircuitBreakerDefaults } from "../StreamFormDefaults"
import { I18n } from "../../../core/i18n/I18n"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"

const FORM_INPUT_STYLE: React.CSSProperties = {
    paddingBottom: 20,
}

const CIRCUIT_BREAKER_SPACING = {
    span: 7,
}

const HIDDEN_BUTTON_PROPS = {
    style: { display: "none" },
}

export interface StreamBackendSettingsModalProps {
    backend: StreamBackend
    validationBasePath: string
    open: boolean
    onClose: () => void
    onChange: (backend: StreamBackend) => void
    validationResult: ValidationResult
    hideWeight?: boolean
}

export default class StreamBackendSettingsModal extends React.Component<StreamBackendSettingsModalProps> {
    private handleChange(attribute: string, value: any) {
        const { backend, onChange } = this.props
        onChange({
            ...backend,
            [attribute]: value,
        })
    }

    render() {
        const { backend, open, onClose, validationResult, validationBasePath, hideWeight } = this.props

        return (
            <Modal
                afterClose={onClose}
                onCancel={onClose}
                title={<I18n id={MessageKey.FrontendStreamComponentsBackendsettingsTitle} />}
                width={800}
                open={open}
                cancelButtonProps={HIDDEN_BUTTON_PROPS}
                okButtonProps={HIDDEN_BUTTON_PROPS}
                closable
                maskClosable
            >
                <Form.Item
                    label={<I18n id={MessageKey.FrontendStreamComponentsBackendsettingsCircuitBreaker} />}
                    required
                >
                    <Form.Item
                        label={<I18n id={MessageKey.CommonEnabled} />}
                        help={<I18n id={MessageKey.FrontendStreamComponentsBackendsettingsEnabledHelp} />}
                        labelCol={CIRCUIT_BREAKER_SPACING}
                        style={FORM_INPUT_STYLE}
                        required
                    >
                        <Switch
                            onChange={value =>
                                this.handleChange("circuitBreaker", value ? streamCircuitBreakerDefaults() : null)
                            }
                            value={backend.circuitBreaker != null}
                        />
                    </Form.Item>
                    <Form.Item
                        label={<I18n id={MessageKey.FrontendStreamComponentsBackendsettingsMaxFailures} />}
                        validateStatus={validationResult.getStatus(`${validationBasePath}.circuitBreaker.maxFailures`)}
                        help={
                            validationResult.getMessage(`${validationBasePath}.circuitBreaker.maxFailures`) ??
                            <I18n id={MessageKey.FrontendStreamComponentsBackendsettingsMaxFailuresHelp} />
                        }
                        labelCol={CIRCUIT_BREAKER_SPACING}
                        style={FORM_INPUT_STYLE}
                        required
                    >
                        <InputNumber
                            min={1}
                            max={9999}
                            style={{ width: "100%" }}
                            onChange={value => this.handleChange("circuitBreaker.maxFailures", value)}
                            value={backend.circuitBreaker?.maxFailures}
                            disabled={backend.circuitBreaker == null}
                        />
                    </Form.Item>
                    <Form.Item
                        label={<I18n id={MessageKey.FrontendStreamComponentsBackendsettingsOpenSeconds} />}
                        validateStatus={validationResult.getStatus(`${validationBasePath}.circuitBreaker.openSeconds`)}
                        help={
                            validationResult.getMessage(`${validationBasePath}.circuitBreaker.openSeconds`) ??
                            <I18n id={MessageKey.FrontendStreamComponentsBackendsettingsOpenSecondsHelp} />
                        }
                        labelCol={CIRCUIT_BREAKER_SPACING}
                        style={hideWeight ? undefined : FORM_INPUT_STYLE}
                        required
                    >
                        <InputNumber
                            min={1}
                            max={9999}
                            style={{ width: "100%" }}
                            onChange={value => this.handleChange("circuitBreaker.openSeconds", value)}
                            value={backend.circuitBreaker?.openSeconds}
                            disabled={backend.circuitBreaker == null}
                        />
                    </Form.Item>
                </Form.Item>
                <Form.Item
                    label={<I18n id={MessageKey.FrontendStreamComponentsBackendsettingsWeight} />}
                    validateStatus={validationResult.getStatus(`${validationBasePath}.weight`)}
                    help={
                        validationResult.getMessage(`${validationBasePath}.weight`) ??
                        <I18n id={MessageKey.FrontendStreamComponentsBackendsettingsWeightHelp} />
                    }
                    hidden={hideWeight}
                >
                    <InputNumber
                        min={0}
                        max={9999}
                        required={false}
                        style={{ width: "100%" }}
                        onChange={value => this.handleChange("weight", value)}
                        value={backend.weight}
                    />
                </Form.Item>
            </Modal>
        )
    }
}
