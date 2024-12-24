import React from "react"
import { Flex, Form, FormInstance, Modal } from "antd"
import { UserOutlined, LogoutOutlined, LockOutlined } from "@ant-design/icons"
import AppContext from "../../../core/components/context/AppContext"
import "./ShellUserMenu.css"
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import UserService from "../UserService"
import { navigateTo } from "../../../core/components/router/AppRouter"
import Notification from "../../../core/components/notification/Notification"
import ValidationResult from "../../../core/validation/ValidationResult"
import Preloader from "../../../core/components/preloader/Preloader"
import UserUpdatePasswordRequest from "../model/UserUpdatePasswordRequest"
import FormLayout from "../../../core/components/form/FormLayout"
import Password from "antd/es/input/Password"
import { UnexpectedResponseError } from "../../../core/apiclient/ApiResponse"
import ValidationResultConverter from "../../../core/validation/ValidationResultConverter"
import { buildLoginUrl } from "../../../core/authentication/buildLoginUrl"

const DEFAULT_FORM_VALUES: UserUpdatePasswordRequest = {
    currentPassword: "",
    newPassword: "",
}

interface ShellUserMenuState {
    modalOpen: boolean
    validationResult: ValidationResult
    loading: boolean
    formValues: UserUpdatePasswordRequest
}

export default class ShellUserMenu extends React.Component<any, ShellUserMenuState> {
    private readonly formRef: React.RefObject<FormInstance | null>
    private readonly service: UserService

    constructor(props: any) {
        super(props)
        this.service = new UserService()
        this.formRef = React.createRef()
        this.state = {
            modalOpen: false,
            validationResult: new ValidationResult(),
            loading: false,
            formValues: DEFAULT_FORM_VALUES,
        }
    }

    private async handleLogout() {
        return UserConfirmation.ask("Are you sure you want to logout?")
            .then(() => this.service.logout())
            .then(() => Notification.success("See ya", "You was logged-out successfully"))
            .then(() => {
                AppContext.get().user = undefined
            })
            .then(() => navigateTo(buildLoginUrl()))
    }

    private async executePasswordChange() {
        const { formValues } = this.state
        this.setState({ validationResult: new ValidationResult() })

        return this.service
            .changePassword(formValues)
            .then(() => Notification.success("Password changed", "Your password was updated successfully"))
            .then(() => this.closeChangePasswordModal())
            .catch(error => this.handleErrorResponse(error))
    }

    private handleErrorResponse(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null) this.setState({ validationResult })
        }

        Notification.error("That didn't work", "Please check the form to see if everything seems correct")
    }

    private openChangePasswordModal() {
        this.setState({ modalOpen: true })
    }

    private closeChangePasswordModal() {
        this.setState({
            modalOpen: false,
            validationResult: new ValidationResult(),
            formValues: DEFAULT_FORM_VALUES,
        })
        this.formRef.current?.resetFields()
    }

    private renderPasswordChangeModal() {
        const { modalOpen, loading, validationResult, formValues } = this.state
        return (
            <Preloader loading={loading}>
                <Modal
                    title="Change password"
                    onCancel={() => this.closeChangePasswordModal()}
                    onClose={() => this.closeChangePasswordModal()}
                    onOk={() => this.executePasswordChange()}
                    open={modalOpen}
                >
                    <Form<UserUpdatePasswordRequest>
                        {...FormLayout.FormDefaults}
                        {...FormLayout.ExpandedLabeledItem}
                        ref={this.formRef}
                        layout="vertical"
                        onValuesChange={(_, formValues) => this.setState({ formValues })}
                        initialValues={formValues}
                    >
                        <Form.Item
                            name="currentPassword"
                            validateStatus={validationResult.getStatus("currentPassword")}
                            help={validationResult.getMessage("currentPassword")}
                            label="Current password"
                            required
                        >
                            <Password />
                        </Form.Item>
                        <Form.Item
                            name="newPassword"
                            validateStatus={validationResult.getStatus("newPassword")}
                            help={validationResult.getMessage("newPassword")}
                            label="New password"
                            required
                        >
                            <Password />
                        </Form.Item>
                    </Form>
                </Modal>
            </Preloader>
        )
    }

    render() {
        const { user } = AppContext.get()
        return (
            <Flex className="shell-user-menu-container">
                <Flex className="shell-user-menu-icon">
                    <UserOutlined />
                </Flex>
                <Flex className="shell-user-menu-user-name">{user?.name}</Flex>
                <Flex className="shell-user-menu-actions">
                    <LockOutlined onClick={() => this.openChangePasswordModal()} />
                    <LogoutOutlined onClick={() => this.handleLogout()} />
                </Flex>

                {this.renderPasswordChangeModal()}
            </Flex>
        )
    }
}
