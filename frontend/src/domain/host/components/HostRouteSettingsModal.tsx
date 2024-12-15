import { Form, FormItemProps, Modal, Switch, Tabs } from "antd"
import FormLayout from "../../../core/components/form/FormLayout"
import TextArea from "antd/es/input/TextArea"
import React from "react"
import ValidationResult from "../../../core/validation/ValidationResult"

const ItemProps: FormItemProps = {
    labelCol: {
        sm: { span: 8 },
    },
    wrapperCol: {
        sm: { span: 14 },
    },
}

export interface HostRouteSettingsProps {
    open: boolean
    index: number
    fieldPath: any
    onClose: () => void
    onCancel: () => void
    validationResult: ValidationResult
}

export default class HostRouteSettingsModal extends React.Component<HostRouteSettingsProps> {
    private renderAdvancedTab() {
        const { index, validationResult, fieldPath } = this.props
        return (
            <>
                <p>
                    Any instruction placed here will be placed in the nginx configuration files as-is. Use this field
                    for any customized configuration parameters that you need in the host route.
                </p>
                <p>
                    Please note that the text below must be in the syntax expected by the nginx. Please refer to the
                    documentation at &nbsp;
                    <a
                        href="https://nginx.org/en/docs/http/ngx_http_core_module.html#location"
                        target="_blank"
                        rel="noreferrer"
                    >
                        this link
                    </a>
                    &nbsp; for more details. If you isn't sure about what to place here, it's probably the best to leave
                    it empty.
                </p>

                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-route-custom-settings"
                    name={[fieldPath, "settings", "custom"]}
                    validateStatus={validationResult.getStatus(`routes[${index}].settings.custom`)}
                    help={validationResult.getMessage(`routes[${index}].settings.custom`)}
                    required
                >
                    <TextArea rows={10} />
                </Form.Item>
            </>
        )
    }

    private renderMainTab() {
        const { index, validationResult, fieldPath } = this.props
        return (
            <>
                <Form.Item
                    {...ItemProps}
                    name={[fieldPath, "settings", "forwardQueryParams"]}
                    label="Forward query params"
                    validateStatus={validationResult.getStatus(`routes[${index}].settings.forwardQueryParams`)}
                    help={
                        validationResult.getMessage(`routes[${index}].settings.forwardQueryParams`) ??
                        "Defines if the query params/string should be forwarded or omitted"
                    }
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    {...ItemProps}
                    name={[fieldPath, "settings", "keepOriginalDomainName"]}
                    label="Keep the original domain name"
                    validateStatus={validationResult.getStatus(`routes[${index}].settings.keepOriginalDomainName`)}
                    help={
                        validationResult.getMessage(`routes[${index}].settings.keepOriginalDomainName`) ??
                        "Defines if the request made by nginx to the target host should use the target's domain as the host"
                    }
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    {...ItemProps}
                    name={[fieldPath, "settings", "proxySslServerName"]}
                    label="Proxy SSL server name"
                    validateStatus={validationResult.getStatus(`routes[${index}].settings.proxySslServerName`)}
                    help={
                        validationResult.getMessage(`routes[${index}].settings.proxySslServerName`) ??
                        "Defines if the SSL negotiation should be made using the target's domain"
                    }
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    {...ItemProps}
                    name={[fieldPath, "settings", "includeForwardHeaders"]}
                    label="Include forward headers"
                    validateStatus={validationResult.getStatus(`routes[${index}].settings.includeForwardHeaders`)}
                    help={
                        validationResult.getMessage(`routes[${index}].settings.includeForwardHeaders`) ??
                        "Defines if headers like 'x-forwarded-for' should be included in the request to the target"
                    }
                    required
                >
                    <Switch />
                </Form.Item>
            </>
        )
    }

    private tabsDefinitions() {
        return [
            {
                key: "main",
                label: "Main",
                children: this.renderMainTab(),
            },
            {
                key: "advanced",
                label: "Advanced",
                children: this.renderAdvancedTab(),
            },
        ]
    }

    render() {
        const { open, onClose, onCancel } = this.props
        return (
            <Modal title="Route settings" open={open} onClose={onClose} onCancel={onCancel} footer={null} width={750}>
                <Tabs defaultActiveKey="1" items={this.tabsDefinitions()} />
            </Modal>
        )
    }
}
