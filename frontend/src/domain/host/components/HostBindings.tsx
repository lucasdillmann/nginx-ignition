import React from "react"
import ValidationResult from "../../../core/validation/ValidationResult"
import { Button, Flex, Form, FormListFieldData, FormListOperation, Input, InputNumber, Select } from "antd"
import If from "../../../core/components/flowcontrol/If"
import { DeleteOutlined, PlusOutlined } from "@ant-design/icons"
import { HostFormBinding } from "../model/HostFormValues"
import { HostBindingType } from "../model/HostRequest"
import PaginatedSelect from "../../../core/components/select/PaginatedSelect"
import CertificateService from "../../certificate/CertificateService"
import PageResponse from "../../../core/pagination/PageResponse"
import { CertificateResponse } from "../../certificate/model/CertificateResponse"
import "./HostBindings.css"
import FormLayout from "../../../core/components/form/FormLayout"
import TagGroup from "../../../core/components/taggroup/TagGroup"
import { I18n } from "../../../core/i18n/I18n"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"

const DEFAULT_VALUES: HostFormBinding = {
    type: HostBindingType.HTTP,
    ip: "0.0.0.0",
    port: 8080,
}

export interface HostBindingsProps {
    pathPrefix: string
    bindings: HostFormBinding[]
    validationResult: ValidationResult
    className?: string
}

export default class HostBindings extends React.Component<HostBindingsProps> {
    private readonly service: CertificateService

    constructor(props: HostBindingsProps) {
        super(props)
        this.service = new CertificateService()
    }

    private async fetchCertificates(
        pageSize: number,
        pageNumber: number,
        searchTerms?: string,
    ): Promise<PageResponse<CertificateResponse>> {
        return this.service.list(pageSize, pageNumber, searchTerms)
    }

    private renderBinding(field: FormListFieldData, operations: FormListOperation, index: number, totalAmount: number) {
        const { validationResult, bindings, pathPrefix } = this.props
        const { name } = field

        return (
            <Flex className="host-form-binding-container">
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-binding-protocol"
                    layout="vertical"
                    name={[name, "type"]}
                    validateStatus={validationResult.getStatus(`${pathPrefix}[${index}].type`)}
                    help={validationResult.getMessage(`${pathPrefix}[${index}].type`)}
                    label={
                        index === 0 ? <I18n id={MessageKey.FrontendHostComponentsHostbindingsProtocol} /> : undefined
                    }
                    required
                >
                    <Select>
                        <Select.Option value={HostBindingType.HTTP}>
                            <I18n id={MessageKey.CommonHttp} />
                        </Select.Option>
                        <Select.Option value={HostBindingType.HTTPS}>
                            <I18n id={MessageKey.CommonHttps} />
                        </Select.Option>
                    </Select>
                </Form.Item>
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-binding-ip"
                    layout="vertical"
                    name={[name, "ip"]}
                    validateStatus={validationResult.getStatus(`${pathPrefix}[${index}].ip`)}
                    help={validationResult.getMessage(`${pathPrefix}[${index}].ip`)}
                    label={
                        index === 0 ? <I18n id={MessageKey.FrontendHostComponentsHostbindingsIpAddress} /> : undefined
                    }
                    required
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-binding-port"
                    layout="vertical"
                    name={[name, "port"]}
                    validateStatus={validationResult.getStatus(`${pathPrefix}[${index}].port`)}
                    help={validationResult.getMessage(`${pathPrefix}[${index}].port`)}
                    label={index === 0 ? <I18n id={MessageKey.CommonPort} /> : undefined}
                    required
                >
                    <InputNumber min={1} max={65535} />
                </Form.Item>

                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-binding-certificate"
                    layout="vertical"
                    name={[name, "certificate"]}
                    validateStatus={validationResult.getStatus(`${pathPrefix}[${index}].certificateId`)}
                    help={validationResult.getMessage(`${pathPrefix}[${index}].certificateId`)}
                    label={
                        index === 0 ? (
                            <I18n id={MessageKey.FrontendHostComponentsHostbindingsSslCertificate} />
                        ) : undefined
                    }
                    required
                >
                    <PaginatedSelect<CertificateResponse>
                        disabled={bindings[index].type === HostBindingType.HTTP}
                        pageProvider={(pageSize, pageNumber, searchTerms) =>
                            this.fetchCertificates(pageSize, pageNumber, searchTerms)
                        }
                        itemKey={certificate => certificate.id}
                        itemDescription={certificate => <TagGroup values={certificate.domainNames} maximumSize={1} />}
                    />
                </Form.Item>
                <If condition={totalAmount > 1}>
                    <DeleteOutlined
                        style={{
                            marginLeft: 15,
                            alignItems: "start",
                            marginTop: index === 0 ? 37 : 7,
                        }}
                        onClick={() => operations.remove(field.name)}
                    />
                </If>
            </Flex>
        )
    }

    private renderBindings(fields: FormListFieldData[], operations: FormListOperation) {
        const bindings = fields.map((field, index) => this.renderBinding(field, operations, index, fields.length))

        const addAction = (
            <Form.Item>
                <Button type="dashed" onClick={() => operations.add(DEFAULT_VALUES)} icon={<PlusOutlined />}>
                    <I18n id={MessageKey.FrontendHostComponentsHostbindingsAddBinding} />
                </Button>
            </Form.Item>
        )

        return [...bindings, addAction]
    }

    render() {
        const { pathPrefix, className } = this.props
        return (
            <div className={className}>
                <Form.List name={pathPrefix}>
                    {(fields, operations) => this.renderBindings(fields, operations)}
                </Form.List>
            </div>
        )
    }
}
