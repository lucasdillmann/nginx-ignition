import React from "react"
import { CertificateResponse } from "./model/CertificateResponse"
import RenewCertificateAction from "./actions/RenewCertificateAction"
import DeleteCertificateAction from "./actions/DeleteCertificateAction"
import { navigateTo, routeParams } from "../../core/components/router/AppRouter"
import CertificateService from "./CertificateService"
import Preloader from "../../core/components/preloader/Preloader"
import AvailableProviderResponse from "./model/AvailableProviderResponse"
import { ProDescriptions } from "@ant-design/pro-components"
import DescriptionLayout from "../../core/components/description/DescriptionLayout"
import { ProFieldValueType } from "@ant-design/pro-utils/es/typing"
import If from "../../core/components/flowcontrol/If"
import "./CertificateDetailsPage.css"
import AppShellContext from "../../core/components/shell/AppShellContext"
import DynamicField, { DynamicFieldType } from "../../core/dynamicfield/DynamicField"
import CommonNotifications from "../../core/components/notification/CommonNotifications"
import EmptyStates from "../../core/components/emptystate/EmptyStates"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import { UserAccessLevel } from "../user/model/UserAccessLevel"
import AccessDeniedPage from "../../core/components/accesscontrol/AccessDeniedPage"

interface CertificateDetailsPageState {
    loading: boolean
    certificate?: CertificateResponse
    availableProviders: AvailableProviderResponse[]
    error?: Error
}

export default class CertificateDetailsPage extends React.Component<unknown, CertificateDetailsPageState> {
    private readonly certificateId: string
    private readonly service: CertificateService

    constructor(props: any) {
        super(props)
        this.certificateId = routeParams().id as string
        this.service = new CertificateService()
        this.state = {
            loading: true,
            availableProviders: [],
        }
    }

    private async deleteCertificate() {
        return DeleteCertificateAction.execute(this.certificateId).then(() => navigateTo("/certificates"))
    }

    private async renewCertificate() {
        return RenewCertificateAction.execute(this.certificateId)
            .then(() => this.setState({ loading: true }))
            .then(() => this.fetchData())
    }

    private updateShellConfig(enableActions: boolean) {
        if (!isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.certificates)) {
            enableActions = false
        }

        AppShellContext.get().updateConfig({
            title: "SSL certificate details",
            subtitle: "Details of a uploaded or issued SSL certificate",
            actions: [
                {
                    description: "Delete",
                    color: "danger",
                    disabled: !enableActions,
                    onClick: () => this.deleteCertificate(),
                },
                {
                    description: "Renew",
                    disabled: !enableActions,
                    onClick: () => this.renewCertificate(),
                },
            ],
        })
    }

    private async fetchData() {
        const certificate = this.service.getById(this.certificateId)
        const providers = this.service.availableProviders()

        return Promise.all([certificate, providers])
            .then(([certificate, availableProviders]) => {
                this.setState({ loading: false, certificate, availableProviders })
                this.updateShellConfig(certificate !== undefined)
            })
            .catch(error => {
                CommonNotifications.failedToFetch()
                this.setState({ loading: false, error })
            })
    }

    private dynamicFieldValue(field: DynamicField) {
        const { certificate } = this.state
        const value = certificate!!.parameters[field.id]

        if (field.type === DynamicFieldType.BOOLEAN) return value ? "Yes / Accepted" : "No / Rejected"

        if (field.type !== DynamicFieldType.ENUM) return value

        return field.enumOptions.find(option => option.id === value)?.description
    }

    private evaluateConditions(field: DynamicField) {
        const { certificate } = this.state
        const { condition } = field
        if (condition === undefined || condition === null) return true

        const { parentField, value } = condition
        return certificate!!.parameters !== undefined && certificate!!.parameters[parentField] === value
    }

    private renderDynamicField(field: DynamicField) {
        if (!this.evaluateConditions(field)) return undefined

        let valueType: ProFieldValueType
        switch (field.type) {
            case DynamicFieldType.EMAIL:
            case DynamicFieldType.FILE:
            case DynamicFieldType.ENUM:
            case DynamicFieldType.SINGLE_LINE_TEXT:
            case DynamicFieldType.BOOLEAN:
            case DynamicFieldType.URL:
                valueType = "text"
                break
            case DynamicFieldType.MULTI_LINE_TEXT:
                valueType = "textarea"
                break
        }

        return (
            <ProDescriptions.Item title={field.description} valueType={valueType}>
                {this.dynamicFieldValue(field)}
            </ProDescriptions.Item>
        )
    }

    private renderDynamicFields(): React.ReactNode[] {
        const { certificate, availableProviders } = this.state
        const provider = availableProviders.find(provider => provider.id === certificate?.providerId)
        const fieldsToRender =
            provider?.dynamicFields
                ?.filter(field => !field.sensitive)
                ?.sort((first, second) => (first.priority > second.priority ? 1 : -1)) ?? []

        return fieldsToRender.map(field => this.renderDynamicField(field))
    }

    private renderContents() {
        const { certificate, availableProviders } = this.state
        const provider = availableProviders.find(provider => provider.id === certificate?.providerId)
        const dateTimeFormat = "DD/MM/YYYY HH:mm:ss"

        return (
            <>
                <h2 className="certificate-details-section-name">General</h2>
                <ProDescriptions {...DescriptionLayout.Defaults} dataSource={certificate}>
                    <ProDescriptions.Item title="Provider">{provider?.name}</ProDescriptions.Item>
                    <ProDescriptions.Item title="Domain names">
                        {certificate?.domainNames.map(domain => (
                            <>
                                {domain}
                                <br />
                            </>
                        ))}
                    </ProDescriptions.Item>
                </ProDescriptions>

                <h2 className="certificate-details-section-name">Validity</h2>
                <ProDescriptions {...DescriptionLayout.Defaults} dataSource={certificate}>
                    <ProDescriptions.Item
                        title="Issued at"
                        dataIndex="issuedAt"
                        valueType="dateTime"
                        fieldProps={{ format: dateTimeFormat }}
                    />
                    <ProDescriptions.Item
                        title="Valid from"
                        dataIndex="validFrom"
                        valueType="dateTime"
                        fieldProps={{ format: dateTimeFormat }}
                    />
                    <ProDescriptions.Item
                        title="Valid until"
                        dataIndex="validUntil"
                        valueType="dateTime"
                        fieldProps={{ format: dateTimeFormat }}
                    />
                    <ProDescriptions.Item
                        title="Renew recommended after"
                        dataIndex="renewAfter"
                        valueType="dateTime"
                        fieldProps={{ format: dateTimeFormat }}
                    />
                </ProDescriptions>

                <If condition={Object.keys(certificate!!.parameters).length > 0}>
                    <h2 className="certificate-details-section-name">Provider-specific parameters</h2>
                    <ProDescriptions {...DescriptionLayout.Defaults} dataSource={certificate}>
                        {this.renderDynamicFields()}
                    </ProDescriptions>
                </If>
            </>
        )
    }

    componentDidMount() {
        this.fetchData()
        this.updateShellConfig(false)
    }

    render() {
        if (!isAccessGranted(UserAccessLevel.READ_ONLY, permissions => permissions.certificates)) {
            return <AccessDeniedPage />
        }

        const { certificate, loading, error } = this.state
        if (loading) return <Preloader loading />
        if (error !== undefined) return EmptyStates.FailedToFetch
        if (certificate === undefined) return EmptyStates.NotFound

        return this.renderContents()
    }
}
