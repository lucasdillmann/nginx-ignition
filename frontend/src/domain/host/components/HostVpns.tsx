import React from "react"
import ValidationResult from "../../../core/validation/ValidationResult"
import { Button, Flex, Form, FormListFieldData, FormListOperation, Input } from "antd"
import { DeleteOutlined, PlusOutlined } from "@ant-design/icons"
import { HostFormVpn } from "../model/HostFormValues"
import PaginatedSelect from "../../../core/components/select/PaginatedSelect"
import FormLayout from "../../../core/components/form/FormLayout"
import VpnResponse from "../../vpn/model/VpnResponse"
import VpnService from "../../vpn/VpnService"
import { vpnRequestDefaults } from "../../vpn/model/VpnRequestDefaults"
import "./HostVpns.css"

export interface HostVpnsProps {
    vpns: HostFormVpn[]
    validationResult: ValidationResult
    className?: string
}

export default class HostVpns extends React.Component<HostVpnsProps> {
    private readonly service: VpnService

    constructor(props: HostVpnsProps) {
        super(props)
        this.service = new VpnService()
    }

    private renderVpn(field: FormListFieldData, operations: FormListOperation, index: number) {
        const { validationResult } = this.props
        const { name } = field

        return (
            <Flex className="host-form-vpn-container">
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-vpn-connection"
                    layout="vertical"
                    name={[name, "vpn"]}
                    validateStatus={validationResult.getStatus(`vpns[${index}].vpnId`)}
                    help={validationResult.getMessage(`vpns[${index}].vpnId`)}
                    label={index === 0 ? "VPN connection" : undefined}
                    required
                >
                    <PaginatedSelect<VpnResponse>
                        itemDescription={item => item.name}
                        itemKey={item => item.id}
                        pageProvider={(pageSize, pageNumber, searchTerms) =>
                            this.service.list(pageSize, pageNumber, searchTerms, true)
                        }
                    />
                </Form.Item>
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-vpn-name"
                    layout="vertical"
                    name={[name, "name"]}
                    validateStatus={validationResult.getStatus(`vpns[${index}].name`)}
                    help={
                        validationResult.getMessage(`vpns[${index}].name`) ??
                        "The unique name to be used in the VPN connection. In Tailscale, this will be used as the " +
                            "subdomain name on your ts.net domain."
                    }
                    label={index === 0 ? "Source name" : undefined}
                    required
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-vpn-host"
                    layout="vertical"
                    name={[name, "host"]}
                    validateStatus={validationResult.getStatus(`vpns[${index}].host`)}
                    help={
                        validationResult.getMessage(`vpns[${index}].host`) ??
                        "The hostname to be used when proxying requests to the target host. If left empty, ignition " +
                            "will use the first option from the domain names above."
                    }
                    label={index === 0 ? "Target host" : undefined}
                    required
                >
                    <Input />
                </Form.Item>

                <DeleteOutlined
                    style={{
                        marginLeft: 15,
                        alignItems: "start",
                        marginTop: index === 0 ? 37 : 7,
                    }}
                    onClick={() => operations.remove(field.name)}
                />
            </Flex>
        )
    }

    private renderVpns(fields: FormListFieldData[], operations: FormListOperation) {
        const vpns = fields.map((field, index) => this.renderVpn(field, operations, index))

        const addAction = (
            <Form.Item>
                <Button type="dashed" onClick={() => operations.add(vpnRequestDefaults())} icon={<PlusOutlined />}>
                    Add VPN binding
                </Button>
            </Form.Item>
        )

        return [...vpns, addAction]
    }

    render() {
        const { className } = this.props
        return (
            <div className={className}>
                <Form.List name="vpns">{(fields, operations) => this.renderVpns(fields, operations)}</Form.List>
            </div>
        )
    }
}
