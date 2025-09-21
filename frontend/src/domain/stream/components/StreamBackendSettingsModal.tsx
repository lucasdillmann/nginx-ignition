import React from "react"
import { Form, InputNumber, Modal, Switch } from "antd"
import ValidationResult from "../../../core/validation/ValidationResult"
import { StreamBackend, StreamCircuitBreaker } from "../model/StreamRequest"

const FORM_INPUT_STYLE: React.CSSProperties = {
    paddingBottom: 20,
}

const CIRCUIT_BREAKER_SPACING = {
    span: 7,
}

const DEFAULT_CIRCUIT_BREAKER: StreamCircuitBreaker = {
    maxFailures: 5,
    openSeconds: 30,
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
                title="Backend settings"
                width={800}
                open={open}
                cancelButtonProps={HIDDEN_BUTTON_PROPS}
                okButtonProps={HIDDEN_BUTTON_PROPS}
                closable
                maskClosable
            >
                <Form.Item label="Circuit breaker" required>
                    <Form.Item
                        label="Enabled"
                        help="Defines if a circuit breaker should be used for this backend or not,
                            automatically interrupting the connection if the threshold is reached."
                        labelCol={CIRCUIT_BREAKER_SPACING}
                        style={FORM_INPUT_STYLE}
                        required
                    >
                        <Switch
                            onChange={value =>
                                this.handleChange("circuitBreaker", value ? { ...DEFAULT_CIRCUIT_BREAKER } : null)
                            }
                            value={backend.circuitBreaker != null}
                        />
                    </Form.Item>
                    <Form.Item
                        label="Maximum failures"
                        validateStatus={validationResult.getStatus(`${validationBasePath}.circuitBreaker.maxFailures`)}
                        help={
                            validationResult.getMessage(`${validationBasePath}.circuitBreaker.maxFailures`) ??
                            "Defines the maximum number of consecutive failures before the circuit breaker is opened"
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
                        label="Open seconds"
                        validateStatus={validationResult.getStatus(`${validationBasePath}.circuitBreaker.openSeconds`)}
                        help={
                            validationResult.getMessage(`${validationBasePath}.circuitBreaker.openSeconds`) ??
                            "Defines amount of time in seconds that the circuit breaker will be open before it is closed again"
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
                    label="Weight"
                    validateStatus={validationResult.getStatus(`${validationBasePath}.weight`)}
                    help={
                        validationResult.getMessage(`${validationBasePath}.weight`) ??
                        "Defines how much of the traffic should be forwarded to this backend among the backend group. " +
                            "The higher the value, the more requests will be routed to it."
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
