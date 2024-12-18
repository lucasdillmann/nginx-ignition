import React from "react"
import { navigateTo, routeParams } from "../../core/components/router/AppRouter"
import AccessListService from "./AccessListService"
import { Form } from "antd"
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
import AccessListRequest, { AccessListOutcome } from "./model/AccessListRequest"
import DeleteAccessListAction from "./actions/DeleteAccessListAction"
import AccessListResponse from "./model/AccessListResponse"
import ReloadNginxAction from "../nginx/actions/ReloadNginxAction"

interface AccessListFormState {
    formValues: AccessListRequest
    validationResult: ValidationResult
    loading: boolean
    notFound: boolean
    error?: Error
}

export default class AccessListFormPage extends React.Component<unknown, AccessListFormState> {
    private readonly accessListId?: string
    private readonly service: AccessListService
    private readonly saveModal: ModalPreloader

    constructor(props: any) {
        super(props)
        const accessListId = routeParams().id
        this.accessListId = accessListId === "new" ? undefined : accessListId
        this.service = new AccessListService()
        this.saveModal = new ModalPreloader()
        this.state = {
            formValues: {
                name: "",
                realm: "",
                defaultOutcome: AccessListOutcome.DENY,
                forwardAuthenticationHeader: false,
                credentials: [],
                entries: [],
            },
            validationResult: new ValidationResult(),
            loading: true,
            notFound: false,
        }
    }

    private submit() {
        const { formValues } = this.state
        this.saveModal.show("Hang on tight", "We're saving the access list")

        const action =
            this.accessListId === undefined
                ? this.service.create(formValues)
                : this.service.updateById(this.accessListId, formValues)

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

    private renderForm() {
        const { formValues } = this.state

        return (
            <Form<AccessListRequest>
                {...FormLayout.FormDefaults}
                onValuesChange={(_, formValues) => this.setState({ formValues })}
                initialValues={formValues}
            >
                <p>TODO: Implement this</p>
            </Form>
        )
    }

    private convertToAccessListRequest(response: AccessListResponse): AccessListRequest {
        const { name, realm, defaultOutcome, forwardAuthenticationHeader, entries, credentials } = response
        return { name, realm, defaultOutcome, forwardAuthenticationHeader, entries, credentials }
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
                    this.setState({ loading: false, formValues: this.convertToAccessListRequest(accessListDetails) })
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
