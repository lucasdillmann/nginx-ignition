import { Form, FormItemProps, Modal, Switch, Tabs } from "antd"
import FormLayout from "../../../core/components/form/FormLayout"
import TextArea from "antd/es/input/TextArea"
import React from "react"
import ValidationResult from "../../../core/validation/ValidationResult"
import { Link } from "react-router-dom"
import PaginatedSelect from "../../../core/components/select/PaginatedSelect"
import AccessListResponse from "../../accesslist/model/AccessListResponse"
import PageResponse from "../../../core/pagination/PageResponse"
import AccessListService from "../../accesslist/AccessListService"

const ItemProps: FormItemProps = {
    labelCol: {
        sm: { span: 8 },
    },
    wrapperCol: {
        sm: { span: 16 },
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
    private accessListService: AccessListService
    constructor(props: HostRouteSettingsProps) {
        super(props)
        this.accessListService = new AccessListService()
    }

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
                    <Link
                        to="https://nginx.org/en/docs/http/ngx_http_core_module.html#location"
                        target="_blank"
                        rel="noreferrer"
                    >
                        this link
                    </Link>
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

    private fetchAccessLists(
        pageSize: number,
        pageNumber: number,
        searchTerms?: string,
    ): Promise<PageResponse<AccessListResponse>> {
        return this.accessListService.list(pageSize, pageNumber, searchTerms)
    }

    private renderMainTab() {
        const { index, validationResult, fieldPath } = this.props
        return (
            <>
                <Form.Item
                    {...ItemProps}
                    name={[fieldPath, "accessList"]}
                    validateStatus={validationResult.getStatus(`routes[${index}].accessListId`)}
                    help={validationResult.getMessage(`routes[${index}].accessListId`)}
                    label="Access list"
                >
                    <PaginatedSelect<AccessListResponse>
                        itemDescription={item => item?.name}
                        itemKey={item => item?.id}
                        pageProvider={(pageSize, pageNumber, searchTerms) =>
                            this.fetchAccessLists(pageSize, pageNumber, searchTerms)
                        }
                        allowEmpty
                    />
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
