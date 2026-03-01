import React from "react"
import ValidationResult from "../../../core/validation/ValidationResult"
import { Button, Flex, Form, FormListFieldData, FormListOperation, Input, Switch } from "antd"
import { DeleteOutlined, PlusOutlined, SettingOutlined } from "@ant-design/icons"
import { HostFormVpn } from "../model/HostFormValues"
import PaginatedSelect from "../../../core/components/select/PaginatedSelect"
import FormLayout from "../../../core/components/form/FormLayout"
import VpnResponse from "../../vpn/model/VpnResponse"
import VpnService from "../../vpn/VpnService"
import { vpnRequestDefaults } from "../../vpn/model/VpnRequestDefaults"
import "./HostVpns.css"
import { I18n } from "../../../core/i18n/I18n"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import EndpointSSLSupport from "../../vpn/model/EndpointSSLSupport"
import { CertificateResponse } from "../../certificate/model/CertificateResponse"
import CertificateService from "../../certificate/CertificateService"
import HostVpnSettingsModal from "./HostVpnSettingsModal"
import If from "../../../core/components/flowcontrol/If"
import TagGroup from "../../../core/components/taggroup/TagGroup"

export interface HostVpnsProps {
    vpns: HostFormVpn[]
    validationResult: ValidationResult
    className?: string
    onChange: () => void
}

interface HostVpnsState {
    vpnSettingsOpenModalIndex?: number
}

export default class HostVpns extends React.Component<HostVpnsProps, HostVpnsState> {
    private readonly service: VpnService
    private readonly certificateService: CertificateService

    constructor(props: HostVpnsProps) {
        super(props)
        this.service = new VpnService()
        this.certificateService = new CertificateService()
        this.state = {}
    }

    private handleVpnChange(index: number, vpn?: VpnResponse) {
        const { vpns, onChange } = this.props
        if (vpn) {
            vpns[index].vpn = vpn
            vpns[index].enableHttps = vpn.driverEndpointSslSupport === EndpointSSLSupport.PROVIDER_MANAGED
            vpns[index].certificate = undefined
            onChange()
        }
    }

    private handleCertificateChange(index: number, certificate?: CertificateResponse) {
        const { vpns, onChange } = this.props
        vpns[index].certificate = certificate
        vpns[index].enableHttps = certificate !== undefined
        onChange()
    }

    private openVpnSettingsModal(index: number) {
        this.setState({ vpnSettingsOpenModalIndex: index })
    }

    private closeVpnSettingsModal() {
        this.setState({ vpnSettingsOpenModalIndex: undefined })
    }

    private renderVpn(field: FormListFieldData, operations: FormListOperation, index: number) {
        const { validationResult, vpns } = this.props
        const { vpnSettingsOpenModalIndex } = this.state
        const { name } = field

        const vpn = vpns[index]?.vpn
        const sslSupport = vpn?.driverEndpointSslSupport

        return (
            <Flex className="host-form-vpn-container">
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-vpn-connection"
                    layout="vertical"
                    name={[name, "vpn"]}
                    validateStatus={validationResult.getStatus(`vpns[${index}].vpnId`)}
                    help={validationResult.getMessage(`vpns[${index}].vpnId`)}
                    label={<I18n id={MessageKey.CommonVpnConnection} />}
                    required
                >
                    <PaginatedSelect<VpnResponse>
                        itemDescription={item => item?.name}
                        itemKey={item => item?.id}
                        pageProvider={(pageSize, pageNumber, searchTerms) =>
                            this.service.list(pageSize, pageNumber, searchTerms, true)
                        }
                        onChange={value => this.handleVpnChange(index, value)}
                    />
                </Form.Item>
                <Form.Item
                    {...FormLayout.ExpandedLabeledItem}
                    className="host-form-vpn-name"
                    layout="vertical"
                    name={[name, "name"]}
                    validateStatus={validationResult.getStatus(`vpns[${index}].name`)}
                    help={
                        validationResult.getMessage(`vpns[${index}].name`) ?? (
                            <I18n id={MessageKey.FrontendHostComponentsHostvpnsPeerNameHelp} />
                        )
                    }
                    label={<I18n id={MessageKey.FrontendHostComponentsHostvpnsPeerName} />}
                    required
                >
                    <Input />
                </Form.Item>

                <If condition={sslSupport === EndpointSSLSupport.PROVIDER_MANAGED}>
                    <Form.Item
                        {...FormLayout.ExpandedLabeledItem}
                        className="host-form-vpn-enable-https"
                        layout="vertical"
                        name={[name, "enableHttps"]}
                        valuePropName="checked"
                        validateStatus={validationResult.getStatus(`vpns[${index}].enableHttps`)}
                        help={
                            validationResult.getMessage(`vpns[${index}].enableHttps`) ?? (
                                <I18n id={MessageKey.FrontendHostComponentsEnableHttpsHelp} />
                            )
                        }
                        label={<I18n id={MessageKey.FrontendHostComponentsEnableHttps} />}
                        required
                    >
                        <Switch />
                    </Form.Item>
                </If>

                <If condition={sslSupport === EndpointSSLSupport.DRIVER_MANAGED}>
                    <Form.Item
                        {...FormLayout.ExpandedLabeledItem}
                        className="host-form-vpn-certificate"
                        layout="vertical"
                        name={[name, "certificate"]}
                        validateStatus={validationResult.getStatus(`vpns[${index}].certificateId`)}
                        help={
                            validationResult.getMessage(`vpns[${index}].certificateId`) ?? (
                                <I18n id={MessageKey.FrontendHostComponentsSslCertificateHelp} />
                            )
                        }
                        label={<I18n id={MessageKey.FrontendHostComponentsSslCertificate} />}
                    >
                        <PaginatedSelect<CertificateResponse>
                            itemDescription={certificate => (
                                <TagGroup values={certificate.domainNames} maximumSize={1} />
                            )}
                            itemKey={item => item?.id}
                            pageProvider={(pageSize, pageNumber, searchTerms) =>
                                this.certificateService.list(pageSize, pageNumber, searchTerms)
                            }
                            onChange={value => this.handleCertificateChange(index, value)}
                            allowEmpty
                        />
                    </Form.Item>
                </If>

                <SettingOutlined
                    style={{
                        marginLeft: 15,
                        alignItems: "start",
                        marginTop: 37,
                        color: "var(--nginxIgnition-colorText)",
                    }}
                    onClick={() => this.openVpnSettingsModal(index)}
                />

                <DeleteOutlined
                    style={{
                        marginLeft: 15,
                        alignItems: "start",
                        marginTop: 37,
                    }}
                    onClick={() => operations.remove(field.name)}
                />

                <HostVpnSettingsModal
                    open={index === vpnSettingsOpenModalIndex}
                    index={index}
                    fieldPath={name}
                    onClose={() => this.closeVpnSettingsModal()}
                    validationResult={validationResult}
                />
            </Flex>
        )
    }

    private renderVpns(fields: FormListFieldData[], operations: FormListOperation) {
        const vpns = fields.map((field, index) => this.renderVpn(field, operations, index))

        const addAction = (
            <Form.Item>
                <Button type="dashed" onClick={() => operations.add(vpnRequestDefaults())} icon={<PlusOutlined />}>
                    <I18n id={MessageKey.FrontendHostComponentsHostvpnsAddBinding} />
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
