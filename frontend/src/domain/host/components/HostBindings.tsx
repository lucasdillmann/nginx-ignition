import React from "react";
import ValidationResult from "../../../core/validation/ValidationResult";
import {Button, Flex, Form, FormListFieldData, FormListOperation, Input, InputNumber, Select} from "antd";
import If from "../../../core/components/flowcontrol/If";
import {CloseOutlined, PlusOutlined} from "@ant-design/icons";
import {HostBinding, HostBindingType} from "../model/HostRequest";
import PaginatedSelect from "../../../core/components/select/PaginatedSelect";
import CertificateService from "../../certificate/CertificateService";
import PageResponse from "../../../core/pagination/PageResponse";
import {CertificateResponse} from "../../certificate/model/CertificateResponse";
import "./HostBindings.css"
import FormLayout from "../../../core/components/form/FormLayout";

const DEFAULT_VALUES: HostBinding = {
    type: HostBindingType.HTTP,
    ip: "0.0.0.0",
    port: 8080,
}

export interface HostBindingsProps {
    bindings: HostBinding[]
    validationResult: ValidationResult
}

export default class HostBindings extends React.Component<HostBindingsProps> {
    private service: CertificateService

    constructor(props: HostBindingsProps) {
        super(props);
        this.service = new CertificateService()
    }

    private async fetchCertificates(pageSize: number, pageNumber: number): Promise<PageResponse<CertificateResponse>> {
        return this.service.list(pageSize, pageNumber)
    }

    private renderBinding(
        field: FormListFieldData,
        operations: FormListOperation,
        index: number,
        totalAmount: number,
    ) {
        const {validationResult, bindings} = this.props
        const {name} = field

        return (
            <Flex className="host-form-binding-container">
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-binding-protocol"
                    layout="vertical"
                    name={[name, "type"]}
                    validateStatus={validationResult.getStatus(`bindings[${index}].type`)}
                    help={validationResult.getMessage(`bindings[${index}].type`)}
                    label={index === 0 ? "Protocol" : undefined}
                    required
                >
                    <Select>
                        <Select.Option value={HostBindingType.HTTP}>HTTP</Select.Option>
                        <Select.Option value={HostBindingType.HTTPS}>HTTPS</Select.Option>
                    </Select>
                </Form.Item>
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-binding-ip"
                    layout="vertical"
                    name={[name, "ip"]}
                    validateStatus={validationResult.getStatus(`bindings[${index}].ip`)}
                    help={validationResult.getMessage(`bindings[${index}].ip`)}
                    label={index === 0 ? "IP address" : undefined}
                    required
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-binding-port"
                    layout="vertical"
                    name={[name, "port"]}
                    validateStatus={validationResult.getStatus(`bindings[${index}].port`)}
                    help={validationResult.getMessage(`bindings[${index}].port`)}
                    label={index === 0 ? "Port" : undefined}
                    required
                >
                    <InputNumber min={1} max={65535} />
                </Form.Item>

                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-binding-certificate"
                    layout="vertical"
                    name={[name, "certificateId"]}
                    validateStatus={validationResult.getStatus(`bindings[${index}].certificateId`)}
                    help={validationResult.getMessage(`bindings[${index}].certificateId`)}
                    label={index === 0 ? "SSL certificate" : undefined}
                    required
                >
                    <PaginatedSelect<CertificateResponse>
                        disabled={bindings[index].type === HostBindingType.HTTP}
                        pageProvider={(pageSize, pageNumber) => this.fetchCertificates(pageSize, pageNumber)}
                        itemKey={certificate => certificate.id}
                        itemDescription={certificate => certificate.domainNames[0]}
                    />
                </Form.Item>
                <If condition={totalAmount > 1}>
                    <CloseOutlined
                        style={{
                            marginLeft: 15,
                            alignItems: "start",
                            marginTop: index === 0 ? 37 : 7
                        }}
                        onClick={() => operations.remove(field.name)}
                    />
                </If>
            </Flex>
        )
    }

    private renderBindings(fields: FormListFieldData[], operations: FormListOperation) {
        const bindings = fields.map(
            (field, index) => this.renderBinding(field, operations, index, fields.length)
        );

        const addAction = (
            <Form.Item>
                <Button
                    type="dashed"
                    onClick={() => operations.add(DEFAULT_VALUES)}
                    icon={<PlusOutlined />}
                >
                    Add binding
                </Button>
            </Form.Item>
        )

        return [...bindings, addAction]
    }

    render() {
        return (
            <Form.List name="bindings">
                {(fields, operations) => this.renderBindings(fields, operations)}
            </Form.List>
        )
    }
}
