import React from "react"
import ValidationResult from "../../core/validation/ValidationResult"
import CertificateService from "./CertificateService"
import AvailableProviderResponse from "./model/AvailableProviderResponse"
import Preloader from "../../core/components/preloader/Preloader"
import { Form, FormInstance, Select } from "antd"
import If from "../../core/components/flowcontrol/If"
import FormLayout from "../../core/components/form/FormLayout"
import DynamicInput from "../../core/components/dynamicfield/DynamicInput"
import { IssueCertificateRequest } from "./model/IssueCertificateRequest"
import ModalPreloader from "../../core/components/preloader/ModalPreloader"
import { IssueCertificateResponse } from "./model/IssueCertificateResponse"
import Notification from "../../core/components/notification/Notification"
import { UnexpectedResponseError } from "../../core/apiclient/ApiResponse"
import ValidationResultConverter from "../../core/validation/ValidationResultConverter"
import DomainNamesList from "../../core/components/domainnames/DomainNamesList"
import { navigateTo } from "../../core/components/router/AppRouter"
import AppShellContext from "../../core/components/shell/AppShellContext"
import { DynamicFieldType } from "../../core/dynamicfield/DynamicField"
import CommonNotifications from "../../core/components/notification/CommonNotifications"
import EmptyStates from "../../core/components/emptystate/EmptyStates"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import { UserAccessLevel } from "../user/model/UserAccessLevel"
import AccessDeniedPage from "../../core/components/accesscontrol/AccessDeniedPage"
import MessageKey from "../../core/i18n/model/MessageKey.generated"
import { raw } from "../../core/i18n/I18n"

interface CertificateIssuePageState {
    availableProviders: AvailableProviderResponse[]
    loading: boolean
    validationResult: ValidationResult
    formValues: IssueCertificateRequest
    error?: Error
}

export default class CertificateIssuePage extends React.Component<unknown, CertificateIssuePageState> {
    private readonly service: CertificateService
    private readonly saveModal: ModalPreloader
    private readonly formRef: React.RefObject<FormInstance | null>

    constructor(props: any) {
        super(props)
        this.service = new CertificateService()
        this.saveModal = new ModalPreloader()
        this.formRef = React.createRef()
        this.state = {
            loading: true,
            validationResult: new ValidationResult(),
            formValues: {
                providerId: "",
                domainNames: [""],
                parameters: {},
            },
            availableProviders: [],
        }
    }

    private async fileToBase64(file: File): Promise<string | null> {
        if (typeof file.arrayBuffer !== "function") return null

        const contents = await file.arrayBuffer()
        return btoa(new Uint8Array(contents).reduce((data, byte) => data + String.fromCharCode(byte), ""))
    }

    private async buildFileParameters() {
        const { availableProviders, formValues } = this.state
        const fileFields =
            availableProviders
                .find(provider => provider.id === formValues.providerId)
                ?.dynamicFields?.filter(field => field.type === DynamicFieldType.FILE) ?? []

        const output: Record<string, any> = {}
        for (const field of fileFields) {
            const value = formValues.parameters?.[field.id]
            if (value !== undefined && value.file !== undefined) output[field.id] = await this.fileToBase64(value.file)
            else output[field.id] = null
        }

        return output
    }

    private async submit() {
        const { formValues } = this.state
        this.saveModal.show(MessageKey.CommonHangOnTight, MessageKey.FrontendCertificateIssuingPleaseWait)
        this.setState({ validationResult: new ValidationResult() })

        const certificateRequest: IssueCertificateRequest = {
            ...formValues,
            domainNames: formValues.domainNames.map(item => item ?? ""),
            parameters: {
                ...formValues.parameters,
                ...(await this.buildFileParameters()),
            },
        }

        this.service
            .issue(certificateRequest)
            .then(response => this.handleResponse(response))
            .catch(error => this.handleError(error))
            .then(() => this.saveModal.close())
    }

    private handleResponse(response: IssueCertificateResponse) {
        const { success, errorReason, certificateId } = response
        if (success) {
            Notification.success(
                MessageKey.FrontendCertificateIssued,
                MessageKey.FrontendCertificateIssueSuccessDescription,
            )
            navigateTo(`/certificates/${certificateId}`)
        } else {
            Notification.error(
                MessageKey.FrontendCertificateIssueFailed,
                errorReason ? raw(errorReason) : MessageKey.FrontendCertificateIssueFailedDescription,
            )
        }
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null) this.setState({ validationResult })
        }

        Notification.error(MessageKey.CommonThatDidntWork, MessageKey.CommonFormCheckMessage)
    }

    private buildProviderSelectOptions() {
        const { availableProviders } = this.state
        return availableProviders.map(provider => ({
            value: provider.id,
            label: provider.name,
        }))
    }

    private fillDynamicFieldsDefaultValues(
        availableProviders: AvailableProviderResponse[],
        formValues: IssueCertificateRequest,
    ): IssueCertificateRequest {
        const { providerId, parameters } = formValues
        const provider = availableProviders.find(item => item.id === providerId)
        if (provider === undefined) return formValues

        const currentParameters = parameters ?? {}
        const updatedParameters = provider.dynamicFields?.reduce(
            (acc, { id, defaultValue }) => ({ [id]: defaultValue, ...acc }),
            currentParameters,
        )

        return { ...formValues, parameters: updatedParameters ?? currentParameters }
    }

    private renderDynamicFields() {
        const { formValues, availableProviders, validationResult } = this.state
        const provider = availableProviders.find(item => item.id === formValues.providerId)
        return provider?.dynamicFields
            .sort((left, right) => (left.priority > right.priority ? 1 : -1))
            .map(field => (
                <DynamicInput
                    key={field.id}
                    validationResult={validationResult}
                    formValues={formValues}
                    field={field}
                />
            ))
    }

    private renderForm() {
        const { validationResult, formValues } = this.state

        return (
            <Form<IssueCertificateRequest>
                {...FormLayout.FormDefaults}
                ref={this.formRef}
                onValuesChange={(changedValues, formValues) => {
                    if (changedValues.providerId !== undefined) {
                        formValues = this.fillDynamicFieldsDefaultValues(this.state.availableProviders, formValues)
                        this.formRef.current?.setFieldsValue(formValues)
                    }
                    this.setState({ formValues })
                }}
                initialValues={formValues}
            >
                <Form.Item
                    name="providerId"
                    validateStatus={validationResult.getStatus("providerId")}
                    help={validationResult.getMessage("providerId")}
                    label="Certificate provider"
                    required
                >
                    <Select placeholder="Certificate provider" options={this.buildProviderSelectOptions()} />
                </Form.Item>
                <DomainNamesList validationResult={validationResult} />
                {this.renderDynamicFields()}
            </Form>
        )
    }

    private updateShellConfig(enableActions: boolean) {
        AppShellContext.get().updateConfig({
            title: MessageKey.FrontendCertificateIssueTitle,
            subtitle: MessageKey.FrontendCertificateIssueSubtitle,
            actions: [
                {
                    description: MessageKey.FrontendCertificateIssueButton,
                    disabled: !enableActions,
                    onClick: () => this.submit(),
                },
            ],
        })
    }

    componentDidMount() {
        this.service
            .availableProviders()
            .then(providers => {
                const sortedProviders = providers.sort((left, right) => (left.priority > right.priority ? 1 : -1))
                const formValues = this.fillDynamicFieldsDefaultValues(sortedProviders, {
                    providerId: providers[0].id,
                    domainNames: [""],
                    parameters: {},
                })
                this.setState({
                    availableProviders: sortedProviders,
                    loading: false,
                    formValues,
                })
                this.formRef.current?.setFieldsValue(formValues)

                this.updateShellConfig(true)
            })
            .catch(error => {
                CommonNotifications.failedToFetch()
                this.setState({ loading: false, error })
            })

        this.updateShellConfig(false)
    }

    render() {
        if (!isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.certificates)) {
            return <AccessDeniedPage />
        }

        const { loading, availableProviders, error } = this.state
        if (error !== undefined) return EmptyStates.FailedToFetch

        return (
            <Preloader loading={loading}>
                <If condition={availableProviders.length > 0}>{this.renderForm()}</If>
            </Preloader>
        )
    }
}
