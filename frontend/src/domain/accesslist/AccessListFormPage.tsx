import React from "react"
import { navigateTo, routeParams } from "../../core/components/router/AppRouter"
import AccessListService from "./AccessListService"
import { Form, FormInstance, Input, Select, Switch } from "antd"
import Preloader from "../../core/components/preloader/Preloader"
import FormLayout from "../../core/components/form/FormLayout"
import ValidationResult from "../../core/validation/ValidationResult"
import ModalPreloader from "../../core/components/preloader/ModalPreloader"
import Notification from "../../core/components/notification/Notification"
import { UnexpectedResponseError } from "../../core/apiclient/ApiResponse"
import ValidationResultConverter from "../../core/validation/ValidationResultConverter"
import AppShellContext, { ShellAction } from "../../core/components/shell/AppShellContext"
import CommonNotifications from "../../core/components/notification/CommonNotifications"
import EmptyStates from "../../core/components/emptystate/EmptyStates"
import DeleteAccessListAction from "./actions/DeleteAccessListAction"
import ReloadNginxAction from "../nginx/actions/ReloadNginxAction"
import AccessListFormDefaults from "./AccessListFormDefaults"
import AccessListFormValues from "./model/AccessListFormValues"
import AccessListConverter from "./AccessListConverter"
import AccessListEntrySets from "./components/AccessListEntrySets"
import AccessListCredentials from "./components/AccessListCredentials"
import { AccessListOutcome } from "./model/AccessListRequest"
import "./AccessListFormPage.css"

interface AccessListFormState {
    formValues: AccessListFormValues
    validationResult: ValidationResult
    loading: boolean
    notFound: boolean
    error?: Error
}

export default class AccessListFormPage extends React.Component<unknown, AccessListFormState> {
    private readonly accessListId?: string
    private readonly service: AccessListService
    private readonly saveModal: ModalPreloader
    private readonly formRef: React.RefObject<FormInstance | null>

    constructor(props: any) {
        super(props)
        const accessListId = routeParams().id
        this.formRef = React.createRef()
        this.accessListId = accessListId === "new" ? undefined : accessListId
        this.service = new AccessListService()
        this.saveModal = new ModalPreloader()
        this.state = {
            formValues: AccessListFormDefaults,
            validationResult: new ValidationResult(),
            loading: true,
            notFound: false,
        }
    }

    private submit() {
        const { formValues } = this.state
        this.saveModal.show("Hang on tight", "We're saving the access list")
        this.setState({ validationResult: new ValidationResult() })

        const data = AccessListConverter.toRequest(formValues)
        const action =
            this.accessListId === undefined
                ? this.service.create(data)
                : this.service.updateById(this.accessListId, data)

        action.then(() => this.handleSuccess()).catch(error => this.handleError(error))
    }

    private handleSuccess() {
        this.saveModal.close()
        navigateTo("/access-lists")
        Notification.success("Access list saved", "The access list was saved successfully")
        ReloadNginxAction.execute()
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null) this.setState({ validationResult })
        }

        this.saveModal.close()
        Notification.error("That didn't work", "Please check the form to see if everything seems correct")
    }

    private removeEntry(index: number) {
        const { formValues } = this.state
        const { entries } = formValues

        let priority = 0
        const updatedValues = entries
            .filter((_, itemIndex) => itemIndex !== index)
            .map(route => ({
                ...route,
                priority: priority++,
            }))

        this.formRef.current?.setFieldValue("entries", updatedValues)
        this.setState({
            formValues: {
                ...formValues,
                entries: updatedValues,
            },
        })
    }

    private handleChange(newValues: AccessListFormValues) {
        const entries = newValues.entries.sort((left, right) => (left.priority > right.priority ? 1 : -1))

        this.setState({
            formValues: {
                ...newValues,
                entries,
            },
        })
    }

    private renderForm() {
        const { formValues, validationResult } = this.state

        return (
            <Form<AccessListFormValues>
                {...FormLayout.FormDefaults}
                ref={this.formRef}
                onValuesChange={(_, formValues) => this.handleChange(formValues)}
                initialValues={formValues}
            >
                <h2 className="access-lists-form-section-name">General</h2>
                <p className="access-lists-form-section-help-text">
                    Main definitions and properties of the access list.
                </p>
                <Form.Item
                    className="access-lists-form-name"
                    name="name"
                    validateStatus={validationResult.getStatus("name")}
                    help={validationResult.getMessage("name")}
                    label="Name"
                    required
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    className="access-lists-form-realm-name"
                    name="realm"
                    validateStatus={validationResult.getStatus("realm")}
                    help={validationResult.getMessage("realm")}
                    label="Realm name"
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    className="access-lists-form-satisfy-all"
                    name="satisfyAll"
                    validateStatus={validationResult.getStatus("satisfyAll")}
                    help={validationResult.getMessage("satisfyAll")}
                    label="Mode"
                    required
                >
                    <Select>
                        <Select.Option value={true}>
                            Both credentials (if any) and entry sets (if any) must be satisfied
                        </Select.Option>
                        <Select.Option value={false}>
                            Either credentials (if any) or entry sets (if any) must be satisfied
                        </Select.Option>
                    </Select>
                </Form.Item>
                <Form.Item
                    className="access-lists-form-default-outcome"
                    name="defaultOutcome"
                    validateStatus={validationResult.getStatus("defaultOutcome")}
                    help={validationResult.getMessage("defaultOutcome")}
                    label="Default outcome"
                    required
                >
                    <Select>
                        <Select.Option value={AccessListOutcome.ALLOW}>Allow access</Select.Option>
                        <Select.Option value={AccessListOutcome.DENY}>Deny access</Select.Option>
                    </Select>
                </Form.Item>
                <Form.Item
                    className="access-lists-form-forward-headers"
                    name="forwardAuthenticationHeader"
                    validateStatus={validationResult.getStatus("forwardAuthenticationHeader")}
                    help={validationResult.getMessage("forwardAuthenticationHeader")}
                    label="Forward authentication headers"
                    required
                >
                    <Switch />
                </Form.Item>

                <h2 className="access-lists-form-section-name">Credentials</h2>
                <p className="access-lists-form-section-help-text">
                    Relation of username and password pairs to be accepted as valid login credentials to access a host
                    or host's route. Leave empty to disable username and password authentication.
                </p>
                <AccessListCredentials validationResult={validationResult} />

                <h2 className="access-lists-form-section-name">Entry sets</h2>
                <p className="access-lists-form-section-help-text">
                    Relation of IP addresses (such as 192.168.0.1) or IP ranges (like 192.168.0.0/24) to either allow or
                    deny the access to the host or host's route. The nginx will evaluate them from top to bottom,
                    executing the first one that matches the source IP address. Leave empty to disable source IP address
                    or range checks.
                </p>
                <AccessListEntrySets
                    entrySets={formValues.entries}
                    validationResult={validationResult}
                    onRemove={index => this.removeEntry(index)}
                />
            </Form>
        )
    }

    private async delete() {
        if (this.accessListId === undefined) return

        return DeleteAccessListAction.execute(this.accessListId).then(() => navigateTo("/access-lists"))
    }

    private updateShellConfig(enableActions: boolean) {
        const actions: ShellAction[] = [
            {
                description: "Save",
                disabled: !enableActions,
                onClick: () => this.submit(),
            },
        ]

        if (this.accessListId !== undefined)
            actions.unshift({
                description: "Delete",
                disabled: !enableActions,
                color: "danger",
                onClick: () => this.delete(),
            })

        AppShellContext.get().updateConfig({
            title: "Access list details",
            subtitle: "Full details and configurations of the access list",
            actions,
        })
    }

    componentDidMount() {
        if (this.accessListId === undefined) {
            this.setState({ loading: false })
            this.updateShellConfig(true)
            return
        }

        this.service
            .getById(this.accessListId!!)
            .then(accessListDetails => {
                if (accessListDetails === undefined) this.setState({ loading: false, notFound: true })
                else {
                    this.setState({ loading: false, formValues: AccessListConverter.toFormValues(accessListDetails) })
                    this.updateShellConfig(true)
                }
            })
            .catch(error => {
                CommonNotifications.failedToFetch()
                this.setState({ loading: false, error })
            })

        this.updateShellConfig(false)
    }

    render() {
        const { loading, notFound, error } = this.state

        if (error !== undefined) return EmptyStates.FailedToFetch
        if (notFound) return EmptyStates.NotFound
        if (loading) return <Preloader loading />

        return this.renderForm()
    }
}
