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
import { BaseOptionType } from "rc-select/lib/Select"
import DynamicInput from "../../core/components/dynamicfield/DynamicInput"

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
        this.saveModal.show("Hang on tight", "We're saving the integration")
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
        Notification.success("Integration saved", "The integration was saved successfully")
        ReloadNginxAction.execute()
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null) this.setState({ validationResult })
        }

        Notification.error("That didn't work", "Please check the form to see if everything seems correct")
    }

    private updateShellConfig(enableActions: boolean) {
        if (!isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.hosts)) {
            enableActions = false
        }

        const actions: ShellAction[] = [
            {
                description: "Save",
                disabled: !enableActions,
                onClick: () => this.submit(),
            },
        ]

        if (this.integrationId !== undefined)
            actions.unshift({
                description: "Delete",
                disabled: !enableActions,
                color: "danger",
                onClick: () => this.delete(),
            })

        AppShellContext.get().updateConfig({
            title: "Integration details",
            subtitle: "Full details and configurations of the integration with third-party apps",
            actions,
        })
    }

    private buildDriverOptions(): BaseOptionType[] {
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
                onValuesChange={(_, formValues) => this.setState({ formValues })}
            >
                <Form.Item
                    name="enabled"
                    validateStatus={validationResult.getStatus("enabled")}
                    help={validationResult.getMessage("enabled")}
                    label="Enabled"
                    required
                >
                    <Switch />
                </Form.Item>
                <Form.Item
                    name="name"
                    validateStatus={validationResult.getStatus("name")}
                    help={validationResult.getMessage("name")}
                    label="Name"
                    required
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    name="driver"
                    validateStatus={validationResult.getStatus("driver")}
                    help={validationResult.getMessage("driver")}
                    label="Driver"
                    required
                >
                    <Select options={this.buildDriverOptions()} />
                </Form.Item>
                {this.renderDynamicFields()}
            </Form>
        )
    }

    componentDidMount() {
        this.updateShellConfig(false)
        const availableDriversPromise = this.service.availableDrivers()

        if (this.integrationId === undefined) {
            return availableDriversPromise.then(availableDrivers => {
                this.setState({ loading: false, availableDrivers })
                this.updateShellConfig(true)
            })
        }

        const dataPromise = this.service.getById(this.integrationId!!)
        Promise.all([dataPromise, availableDriversPromise])
            .then(([formValues, availableDrivers]) => {
                const notFound = formValues === undefined

                this.setState({
                    loading: false,
                    formValues: formValues ?? integrationRequestDefaults(),
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
