import React from "react"
import { Card, Empty, Form } from "antd"
import ValidationResult from "../../../core/validation/ValidationResult"
import { HostBinding, HostBindingType } from "../model/HostRequest"
import PaginatedSelect from "../../../core/components/select/PaginatedSelect"
import CertificateService from "../../certificate/CertificateService"
import PageResponse from "../../../core/pagination/PageResponse"
import { CertificateResponse } from "../../certificate/model/CertificateResponse"
import TagGroup from "../../../core/components/taggroup/TagGroup"

export interface HostGlobalBindingCertificateOverridesProps {
    globalBindings?: HostBinding[]
    validationResult: ValidationResult
}

export default class HostGlobalBindingCertificateOverrides extends React.Component<HostGlobalBindingCertificateOverridesProps> {
    private readonly certificateService: CertificateService

    constructor(props: HostGlobalBindingCertificateOverridesProps) {
        super(props)
        this.certificateService = new CertificateService()
    }

    private getBindingLabel(binding: HostBinding): string {
        const protocol = binding.type === HostBindingType.HTTPS ? "HTTPS" : "HTTP"
        return `${protocol} - ${binding.ip}:${binding.port}`
    }

    private fetchCertificates = async (search: string, page: number): Promise<PageResponse<CertificateResponse>> => {
        return this.certificateService.list(undefined, page, search)
    }

    private renderBinding(binding: HostBinding, index: number) {
        const { validationResult } = this.props
        const fieldName = ["globalBindingCertificateOverrides", index, "certificate"]

        if (binding.type !== HostBindingType.HTTPS) {
            return null
        }

        return (
            <Card key={binding.id || `binding-${index}`} style={{ marginBottom: 16 }} size="small">
                <Form.Item
                    name={fieldName}
                    label={this.getBindingLabel(binding)}
                    validateStatus={validationResult.getStatus(
                        `globalBindingCertificateOverrides[${index}].certificate`,
                    )}
                    help={validationResult.getMessage(`globalBindingCertificateOverrides[${index}].certificate`)}
                    labelCol={{ span: 8 }}
                    wrapperCol={{ span: 16 }}
                >
                    <PaginatedSelect<CertificateResponse>
                        placeholder="Use default certificate"
                        allowEmpty={true}
                        pageProvider={(pageSize, pageNumber, searchTerms) =>
                            this.fetchCertificates(searchTerms || "", pageNumber)
                        }
                        itemKey={certificate => certificate.id}
                        itemDescription={certificate => <TagGroup values={certificate.domainNames} maximumSize={1} />}
                    />
                </Form.Item>
                <Form.Item
                    name={["globalBindingCertificateOverrides", index, "bindingId"]}
                    hidden
                    initialValue={binding.id}
                >
                    <input type="hidden" value={binding.id} />
                </Form.Item>
            </Card>
        )
    }

    render() {
        const { globalBindings } = this.props

        if (!globalBindings || globalBindings.length === 0) {
            return (
                <Empty
                    description="No global bindings configured. Please configure global bindings in the settings."
                    style={{ marginTop: 20, marginBottom: 20 }}
                />
            )
        }

        const httpsBindings = globalBindings.filter(b => b.type === HostBindingType.HTTPS)

        if (httpsBindings.length === 0) {
            return (
                <Empty
                    description="No HTTPS global bindings configured. Certificate overrides are only available for HTTPS bindings."
                    style={{ marginTop: 20, marginBottom: 20 }}
                />
            )
        }

        return (
            <div style={{ marginTop: 16 }}>
                <p style={{ marginBottom: 16, color: "#666" }}>
                    Select a certificate for each HTTPS binding to override the default certificate. Leave empty to use
                    the certificate configured in the global binding.
                </p>
                {globalBindings.map((binding, index) => this.renderBinding(binding, index))}
            </div>
        )
    }
}
