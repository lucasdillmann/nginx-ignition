import React from "react"
import { navigateTo, routeParams } from "../../core/components/router/AppRouter"
import UserRequest from "./model/UserRequest"
import UserService from "./UserService"
import { Form, Input, Switch } from "antd"
import Preloader from "../../core/components/preloader/Preloader"
import FormLayout from "../../core/components/form/FormLayout"
import ValidationResult from "../../core/validation/ValidationResult"
import Password from "antd/es/input/Password"
import ModalPreloader from "../../core/components/preloader/ModalPreloader"
import Notification from "../../core/components/notification/Notification"
import { UnexpectedResponseError } from "../../core/apiclient/ApiResponse"
import ValidationResultConverter from "../../core/validation/ValidationResultConverter"
import UserResponse from "./model/UserResponse"
import AppShellContext, { ShellAction } from "../../core/components/shell/AppShellContext"
import DeleteUserAction from "./actions/DeleteUserAction"
import MessageKey from "../../core/i18n/model/MessageKey.generated"
import AppContext from "../../core/components/context/AppContext"
import CommonNotifications from "../../core/components/notification/CommonNotifications"
import EmptyStates from "../../core/components/emptystate/EmptyStates"
import { UserPermissionToggle } from "./components/UserPermissionToggle"
import { UserAccessLevel } from "./model/UserAccessLevel"
import AccessDeniedPage from "../../core/components/accesscontrol/AccessDeniedPage"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"

interface UserFormState {
    formValues: UserRequest
    validationResult: ValidationResult
    loading: boolean
    notFound: boolean
    error?: Error
}

export default class UserFormPage extends React.Component<unknown, UserFormState> {
    private readonly service: UserService
    private readonly saveModal: ModalPreloader
    private userId?: string

    constructor(props: any) {
        super(props)
        const userId = routeParams().id
        this.userId = userId === "new" ? undefined : userId
        this.service = new UserService()
        this.saveModal = new ModalPreloader()
        this.state = {
            formValues: {
                name: "",
                enabled: true,
                username: "",
                permissions: {
                    hosts: UserAccessLevel.READ_WRITE,
                    streams: UserAccessLevel.READ_WRITE,
                    certificates: UserAccessLevel.READ_WRITE,
                    logs: UserAccessLevel.READ_ONLY,
                    integrations: UserAccessLevel.READ_WRITE,
                    accessLists: UserAccessLevel.READ_WRITE,
                    settings: UserAccessLevel.READ_WRITE,
                    users: UserAccessLevel.READ_WRITE,
                    nginxServer: UserAccessLevel.READ_WRITE,
                    exportData: UserAccessLevel.READ_ONLY,
                    vpns: UserAccessLevel.READ_WRITE,
                    caches: UserAccessLevel.READ_WRITE,
                },
            },
            validationResult: new ValidationResult(),
            loading: true,
            notFound: false,
        }
    }

    private submit() {
        const { formValues } = this.state
        this.saveModal.show(MessageKey.CommonHangOnTight, {
            id: MessageKey.CommonSavingType,
            params: { type: MessageKey.CommonEntityUser },
        })
        this.setState({ validationResult: new ValidationResult() })

        const action =
            this.userId === undefined
                ? this.service.create(formValues).then(response => this.updateId(response.id))
                : this.service.updateById(this.userId, formValues)

        action
            .then(() => this.handleSuccess())
            .catch(error => this.handleError(error))
            .then(() => this.saveModal.close())
    }

    private updateId(id: string) {
        this.userId = id
        navigateTo(`/users/${id}`, true)
        this.updateShellConfig(true)
    }

    private handleSuccess() {
        Notification.success(
            { id: MessageKey.CommonTypeSaved, params: { type: MessageKey.CommonEntityUser } },
            MessageKey.CommonSuccessMessage,
        )
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null) this.setState({ validationResult })
        }

        Notification.error(MessageKey.CommonThatDidntWork, MessageKey.CommonFormCheckMessage)
    }

    private passwordHelpText() {
        return this.userId === undefined ? undefined : "Leave empty if you want to keep the user's password unchanged"
    }

    private renderForm() {
        const { validationResult, formValues } = this.state

        return (
            <Form<UserRequest>
                {...FormLayout.FormDefaults}
                onValuesChange={(_, formValues) => this.setState({ formValues })}
                initialValues={formValues}
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
                    name="username"
                    validateStatus={validationResult.getStatus("username")}
                    help={validationResult.getMessage("username")}
                    label="Username"
                    required
                >
                    <Input />
                </Form.Item>
                <Form.Item
                    name="password"
                    validateStatus={validationResult.getStatus("password")}
                    help={validationResult.getMessage("password") ?? this.passwordHelpText()}
                    label="Password"
                    required={this.userId === undefined}
                >
                    <Password />
                </Form.Item>
                <Form.Item label="Permissions" required>
                    <UserPermissionToggle id="hosts" label="Hosts" />
                    <UserPermissionToggle id="streams" label="Streams" />
                    <UserPermissionToggle id="certificates" label="SSL certificates" />
                    <UserPermissionToggle id="integrations" label="Integrations" />
                    <UserPermissionToggle id="vpns" label="VPNs" />
                    <UserPermissionToggle id="caches" label="Cache configurations" />
                    <UserPermissionToggle id="accessLists" label="Access lists" />
                    <UserPermissionToggle id="settings" label="Settings" />
                    <UserPermissionToggle id="users" label="Users" />
                    <UserPermissionToggle id="logs" label="Logs" disableReadWrite />
                    <UserPermissionToggle id="exportData" label="Export and backup" disableReadWrite />
                    <UserPermissionToggle id="nginxServer" label="Nginx server control" disableNoAccess />
                </Form.Item>
            </Form>
        )
    }

    private convertToUserRequest(response: UserResponse): UserRequest {
        const { enabled, name, username, permissions } = response
        return { enabled, name, username, permissions }
    }

    private async delete() {
        if (this.userId === undefined) return

        return DeleteUserAction.execute(this.userId).then(() => navigateTo("/users"))
    }

    private updateShellConfig(enableActions: boolean) {
        if (!isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.users)) {
            enableActions = false
        }

        const actions: ShellAction[] = [
            {
                description: MessageKey.CommonSave,
                disabled: !enableActions,
                onClick: () => this.submit(),
            },
        ]

        if (this.userId !== undefined)
            actions.unshift({
                description: MessageKey.CommonDelete,
                disabled: !enableActions || this.userId === AppContext.get().user?.id,
                color: "danger",
                onClick: () => this.delete(),
            })

        AppShellContext.get().updateConfig({
            title: MessageKey.FrontendUserFormTitle,
            subtitle: MessageKey.FrontendUserFormSubtitle,
            actions,
        })
    }

    componentDidMount() {
        if (this.userId === undefined) {
            this.setState({ loading: false })
            this.updateShellConfig(true)
            return
        }

        this.service
            .getById(this.userId!!)
            .then(userDetails => {
                if (userDetails === undefined) this.setState({ loading: false, notFound: true })
                else {
                    this.setState({ loading: false, formValues: this.convertToUserRequest(userDetails) })
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
        if (!isAccessGranted(UserAccessLevel.READ_ONLY, permissions => permissions.users)) {
            return <AccessDeniedPage />
        }

        const { loading, notFound, error } = this.state

        if (error !== undefined) return EmptyStates.FailedToFetch
        if (notFound) return EmptyStates.NotFound
        if (loading) return <Preloader loading />

        return this.renderForm()
    }
}
