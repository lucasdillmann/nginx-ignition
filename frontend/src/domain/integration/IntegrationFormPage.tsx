import React from "react"
import AppShellContext, { ShellAction } from "../../core/components/shell/AppShellContext"
import { Form, FormInstance, Input, Select, Switch } from "antd"
import FormLayout from "../../core/components/form/FormLayout"
import { navigateTo, routeParams } from "../../core/components/router/AppRouter"
import ValidationResult from "../../core/validation/ValidationResult"
import Preloader from "../../core/components/preloader/Preloader"
import Notification from "../../core/components/notification/Notification"
import { UnexpectedResponseError } from "../../core/apiclient/ApiResponse"
import ValidationResultConverter from "../../core/validation/ValidationResultConverter"
import ModalPreloader from "../../core/components/preloader/ModalPreloader"
import ReloadNginxAction from "../nginx/actions/ReloadNginxAction"
import CommonNotifications from "../../core/components/notification/CommonNotifications"
import EmptyStates from "../../core/components/emptystate/EmptyStates"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import { UserAccessLevel } from "../user/model/UserAccessLevel"
import AccessDeniedPage from "../../core/components/accesscontrol/AccessDeniedPage"
import IntegrationService from "./IntegrationService"
import IntegrationRequest from "./model/IntegrationRequest"
import DeleteIntegrationAction from "./actions/DeleteIntegrationAction"
import { integrationRequestDefaults } from "./model/IntegrationRequestDefaults"
import AvailableDriverResponse from "./model/AvailableDriverResponse"
import DynamicInput from "../../core/components/dynamicfield/DynamicInput"
import MessageKey from "../../core/i18n/model/MessageKey.generated"
import { I18n } from "../../core/i18n/I18n"

interface IntegrationFormPageState {
    availableDrivers: AvailableDriverResponse[]
    formValues: IntegrationRequest
    validationResult: ValidationResult
    loading: boolean
    notFound: boolean
    error?: Error
}

export default class IntegrationFormPage extends React.Component<any, IntegrationFormPageState> {
    private readonly service: IntegrationService
    private readonly saveModal: ModalPreloader
    private readonly formRef: React.RefObject<FormInstance | null>
    private integrationId?: string

    constructor(props: any) {
        super(props)

        const hostId = routeParams().id
        this.integrationId = hostId === "new" ? undefined : hostId
        this.service = new IntegrationService()
        this.saveModal = new ModalPreloader()
        this.formRef = React.createRef()
        this.state = {
            availableDrivers: [],
            validationResult: new ValidationResult(),
            loading: true,
            notFound: false,
            formValues: integrationRequestDefaults(),
        }
    }

    private async delete() {
        if (this.integrationId === undefined) return

        return DeleteIntegrationAction.execute(this.integrationId).then(() => navigateTo("/integrations"))
    }

    private submit() {
        const { formValues } = this.state
        this.saveModal.show(MessageKey.CommonHangOnTight, {
            id: MessageKey.CommonSavingType,
            params: { type: MessageKey.CommonIntegration },
        })
        this.setState({ validationResult: new ValidationResult() })

        const action =
            this.integrationId === undefined
                ? this.service.create(formValues).then(response => this.updateId(response.id))
                : this.service.updateById(this.integrationId, formValues)

        action
            .then(() => this.handleSuccess())
            .catch(error => this.handleError(error))
            .then(() => this.saveModal.close())
    }

    private updateId(id: string) {
        this.integrationId = id
        navigateTo(`/integrations/${id}`, true)
        this.updateShellConfig(true)
    }

    private handleSuccess() {
        Notification.success(
            { id: MessageKey.CommonTypeSaved, params: { type: MessageKey.CommonIntegration } },
            MessageKey.CommonSuccessMessage,
        )
        ReloadNginxAction.execute()
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null) this.setState({ validationResult })
        }

        Notification.error(MessageKey.CommonThatDidntWork, MessageKey.CommonFormCheckMessage)
    }

    private updateShellConfig(enableActions: boolean) {
        if (!isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.hosts)) {
            enableActions = false
        }

        const actions: ShellAction[] = [
            {
                description: MessageKey.CommonSave,
                disabled: !enableActions,
                onClick: () => this.submit(),
            },
        ]

        if (this.integrationId !== undefined)
            actions.unshift({
                description: MessageKey.CommonDelete,
                disabled: !enableActions,
                color: "danger",
                onClick: () => this.delete(),
            })

        AppShellContext.get().updateConfig({
            title: MessageKey.FrontendIntegrationFormTitle,
            subtitle: MessageKey.FrontendIntegrationFormSubtitle,
            actions,
        })
    }

    private buildDriverOptions(): any[] {
        const { availableDrivers } = this.state
        return availableDrivers.map(driver => ({
            value: driver.id,
            label: driver.name,
        }))
    }

    private renderDynamicFields() {
        const { formValues, availableDrivers, validationResult } = this.state
        const provider = availableDrivers.find(item => item.id === formValues.driver)
        return provider?.configurationFields
            .sort((left, right) => (left.priority > right.priority ? 1 : -1))
            .map(field => (
                <DynamicInput
                    key={field.id}
                    dataField="parameters"
                    validationResult={validationResult}
                    formValues={formValues}
                    field={field}
                />
            ))
    }

    private renderForm() {
        const { validationResult, formValues } = this.state

        return (
            <Form<IntegrationRequest>
                {...FormLayout.FormDefaults}
                ref={this.formRef}
                initialValues={formValues}
                onValuesChange={(changedValues, formValues) => {
                    if (changedValues.driver !== undefined) {
                        formValues = this.fillDynamicFieldsDefaultValues(formValues, this.state.availableDrivers)
                        this.formRef.current?.setFieldsValue(formValues)
                    }
                    this.setState({ formValues })
                }}
            >
                <Form.Item
                    name="enabled"
                    validateStatus={validationResult.getStatus("enabled")}
                    help={validationResult.getMessage("enabled")}
                    label={<I18n id={MessageKey.CommonEnabled} />}
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name="name"
                    validateStatus={validationResult.getStatus("name")}
                    help={validationResult.getMessage("name")}
                    label={<I18n id={MessageKey.CommonName} />}
                    required
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    name="driver"
                    validateStatus={validationResult.getStatus("driver")}
                    help={validationResult.getMessage("driver")}
                    label={<I18n id={MessageKey.CommonDriver} />}
                    required
                >
                    <Select options={this.buildDriverOptions()} />
                </Form.Item>
                {this.renderDynamicFields()}
            </Form>
        )
    }

    private fillDynamicFieldsDefaultValues(
        formValues: IntegrationRequest,
        availableDrivers: AvailableDriverResponse[],
    ): IntegrationRequest {
        const { driver, parameters } = formValues
        const provider = availableDrivers.find(item => item.id === driver)
        if (provider === undefined) return formValues

        const currentParameters = parameters ?? {}
        const updatedParameters = provider.configurationFields?.reduce(
            (acc, { id, defaultValue }) => ({ [id]: defaultValue, ...acc }),
            currentParameters,
        )

        return { ...formValues, parameters: updatedParameters ?? currentParameters }
    }

    componentDidMount() {
        this.updateShellConfig(false)
        const availableDriversPromise = this.service.availableDrivers()

        if (this.integrationId === undefined) {
            return availableDriversPromise.then(availableDrivers => {
                const formValues = this.fillDynamicFieldsDefaultValues(
                    { ...integrationRequestDefaults(), driver: availableDrivers[0]?.id },
                    availableDrivers,
                )
                this.setState({ loading: false, availableDrivers, formValues })
                this.formRef.current?.setFieldsValue(formValues)
                this.updateShellConfig(true)
            })
        }

        const dataPromise = this.service.getById(this.integrationId!!)
        Promise.all([dataPromise, availableDriversPromise])
            .then(([formValues, availableDrivers]) => {
                const notFound = formValues === undefined

                this.setState({
                    loading: false,
                    formValues: this.fillDynamicFieldsDefaultValues(
                        formValues ?? integrationRequestDefaults(),
                        availableDrivers,
                    ),
                    notFound,
                    availableDrivers,
                })
                this.updateShellConfig(!notFound)
            })
            .catch(error => {
                CommonNotifications.failedToFetch()
                this.setState({ loading: false, error })
            })
    }

    render() {
        if (!isAccessGranted(UserAccessLevel.READ_ONLY, permissions => permissions.integrations)) {
            return <AccessDeniedPage />
        }

        const { loading, notFound, error } = this.state

        if (error !== undefined) return EmptyStates.FailedToFetch
        if (notFound) return EmptyStates.NotFound
        if (loading) return <Preloader loading />

        return this.renderForm()
    }
}
