import React from "react"
import { Button, Flex, Form, FormInstance, Input, Modal, Tabs, Typography } from "antd"
import { LockOutlined, SafetyOutlined, UserOutlined } from "@ant-design/icons"
import UserService from "../UserService"
import Notification from "../../../core/components/notification/Notification"
import ValidationResult from "../../../core/validation/ValidationResult"
import Preloader from "../../../core/components/preloader/Preloader"
import UserUpdatePasswordRequest from "../model/UserUpdatePasswordRequest"
import UserUpdateProfileRequest from "../model/UserUpdateProfileRequest"
import FormLayout from "../../../core/components/form/FormLayout"
import Password from "antd/es/input/Password"
import { UnexpectedResponseError } from "../../../core/apiclient/ApiResponse"
import ValidationResultConverter from "../../../core/validation/ValidationResultConverter"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { I18n } from "../../../core/i18n/I18n"
import TotpSetup from "./TotpSetup"
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import AppContext from "../../../core/components/context/AppContext"
import "./UserSecuritySettingsModal.css"

export type UserSecuritySettingsTab = "profile" | "password" | "totp"

interface UserSecuritySettingsModalProps {
    open: boolean
    onCancel: () => void
    initialTab?: UserSecuritySettingsTab
}

interface UserSecuritySettingsModalState {
    loading: boolean
    validationResult: ValidationResult
    passwordFormValues: UserUpdatePasswordRequest
    profileFormValues: UserUpdateProfileRequest
    totpEnabled?: boolean
    totpLoading: boolean
}

const DEFAULT_PASSWORD_FORM_VALUES: UserUpdatePasswordRequest = {
    currentPassword: "",
    newPassword: "",
}

const DEFAULT_PROFILE_FORM_VALUES: UserUpdateProfileRequest = {
    name: "",
    username: "",
}

export default class UserSecuritySettingsModal extends React.Component<
    UserSecuritySettingsModalProps,
    UserSecuritySettingsModalState
> {
    private readonly passwordFormRef: React.RefObject<FormInstance | null>
    private readonly profileFormRef: React.RefObject<FormInstance | null>
    private readonly service: UserService

    constructor(props: UserSecuritySettingsModalProps) {
        super(props)
        this.service = new UserService()
        this.passwordFormRef = React.createRef()
        this.profileFormRef = React.createRef()
        this.state = {
            loading: false,
            validationResult: new ValidationResult(),
            passwordFormValues: DEFAULT_PASSWORD_FORM_VALUES,
            profileFormValues: DEFAULT_PROFILE_FORM_VALUES,
            totpLoading: true,
        }
    }

    componentDidUpdate(prevProps: Readonly<UserSecuritySettingsModalProps>) {
        if (this.props.open && !prevProps.open) {
            this.fetchTotpStatus()
            this.loadProfileFormValues()
        }
    }

    private loadProfileFormValues() {
        const user = AppContext.get().user
        const profileFormValues = {
            name: user?.name ?? "",
            username: user?.username ?? "",
        }

        this.setState({
            profileFormValues,
            validationResult: new ValidationResult(),
        })
        this.profileFormRef.current?.setFieldsValue(profileFormValues)
    }

    private fetchTotpStatus() {
        this.setState({ totpLoading: true })
        this.service
            .getTotpStatus()
            .then(enabled => this.setState({ totpEnabled: enabled, totpLoading: false }))
            .catch(() => this.setState({ totpLoading: false }))
    }

    private async executeProfileUpdate() {
        const { profileFormValues } = this.state
        this.setState({ validationResult: new ValidationResult(), loading: true })

        return this.service
            .updateProfile(profileFormValues)
            .then(() =>
                Notification.success(
                    { id: MessageKey.CommonTypeSaved, params: { type: MessageKey.CommonUser } },
                    MessageKey.CommonSuccessMessage,
                ),
            )
            .then(() => AppContext.get().container!!.reload())
            .then(() => this.props.onCancel())
            .catch(error => this.handleErrorResponse(error))
            .then(() => this.setState({ loading: false }))
    }

    private async executePasswordChange() {
        const { passwordFormValues } = this.state
        this.setState({ validationResult: new ValidationResult(), loading: true })

        return this.service
            .changePassword(passwordFormValues)
            .then(() => Notification.success(MessageKey.CommonPasswordChanged, MessageKey.CommonSuccessMessage))
            .then(() => this.resetPasswordForm())
            .catch(error => this.handleErrorResponse(error))
            .then(() => this.setState({ loading: false }))
    }

    private handleErrorResponse(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null) this.setState({ validationResult })
        }

        Notification.error(MessageKey.CommonThatDidntWork, MessageKey.CommonFormCheckMessage)
    }

    private resetPasswordForm() {
        this.setState({
            validationResult: new ValidationResult(),
            passwordFormValues: DEFAULT_PASSWORD_FORM_VALUES,
        })
        this.passwordFormRef.current?.resetFields()
        this.props.onCancel()
    }

    private handleTotpDisabled() {
        UserConfirmation.ask(MessageKey.FrontendUserMenuTotpDisableConfirmation)
            .then(() => this.setState({ totpLoading: true }))
            .then(() => this.service.disableTotp())
            .then(() =>
                Notification.success(
                    MessageKey.FrontendUserMenuTotpDisabledTitle,
                    MessageKey.FrontendUserMenuTotpDisabledSuccessDescription,
                ),
            )
            .catch(() => Notification.error(MessageKey.CommonThatDidntWork, MessageKey.CommonTryAgainLater))
            .then(() => this.setState({ totpLoading: false }))
    }

    private renderProfileTab() {
        const { validationResult, profileFormValues } = this.state
        return (
            <Form<UserUpdateProfileRequest>
                {...FormLayout.FormDefaults}
                labelCol={FormLayout.ExpandedLabeledItem.labelCol}
                wrapperCol={FormLayout.ExpandedLabeledItem.wrapperCol}
                ref={this.profileFormRef}
                layout="vertical"
                onValuesChange={(_, profileFormValues) => this.setState({ profileFormValues })}
                initialValues={profileFormValues}
                className="user-security-settings-form"
            >
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
                    name="username"
                    validateStatus={validationResult.getStatus("username")}
                    help={validationResult.getMessage("username")}
                    label={<I18n id={MessageKey.CommonUsername} />}
                    required
                >
                    <Input />
                </Form.Item>

                <Flex justify="end" style={{ marginTop: 24 }}>
                    <Button type="primary" onClick={() => this.executeProfileUpdate()}>
                        <I18n id={MessageKey.CommonSave} />
                    </Button>
                </Flex>
            </Form>
        )
    }

    private renderPasswordTab() {
        const { validationResult, passwordFormValues } = this.state
        return (
            <Form<UserUpdatePasswordRequest>
                {...FormLayout.FormDefaults}
                labelCol={FormLayout.ExpandedLabeledItem.labelCol}
                wrapperCol={FormLayout.ExpandedLabeledItem.wrapperCol}
                ref={this.passwordFormRef}
                layout="vertical"
                onValuesChange={(_, passwordFormValues) => this.setState({ passwordFormValues })}
                initialValues={passwordFormValues}
                className="user-security-settings-form"
            >
                <Form.Item
                    name="currentPassword"
                    validateStatus={validationResult.getStatus("currentPassword")}
                    help={validationResult.getMessage("currentPassword")}
                    label={<I18n id={MessageKey.FrontendUserMenuCurrentPassword} />}
                    required
                >
                    <Password />
                </Form.Item>
                <Form.Item
                    name="newPassword"
                    validateStatus={validationResult.getStatus("newPassword")}
                    help={validationResult.getMessage("newPassword")}
                    label={<I18n id={MessageKey.FrontendUserMenuNewPassword} />}
                    required
                >
                    <Password />
                </Form.Item>

                <Flex justify="end" style={{ marginTop: 24 }}>
                    <Button type="primary" onClick={() => this.executePasswordChange()}>
                        <I18n id={MessageKey.CommonSave} />
                    </Button>
                </Flex>
            </Form>
        )
    }

    private renderTotpTab() {
        const { totpEnabled, totpLoading } = this.state

        if (totpLoading) {
            return <Preloader loading={true} />
        }

        if (totpEnabled) {
            return (
                <Flex vertical align="center" justify="center" className="totp-enabled-container">
                    <SafetyOutlined className="totp-enabled-icon" />
                    <Typography.Title level={4}>
                        <I18n id={MessageKey.FrontendUserMenuTotpEnabledTitle} />
                    </Typography.Title>
                    <Typography.Text type="secondary" className="totp-enabled-description">
                        <I18n id={MessageKey.FrontendUserMenuTotpEnabledDescription} />
                    </Typography.Text>
                    <Button
                        danger
                        type="primary"
                        onClick={() => this.handleTotpDisabled()}
                        className="totp-disable-button"
                    >
                        <I18n id={MessageKey.FrontendUserMenuTotpDisableButton} />
                    </Button>
                </Flex>
            )
        }

        return <TotpSetup onActivation={() => this.setState({ totpEnabled: true })} />
    }

    render() {
        const { open, onCancel, initialTab = "password" } = this.props
        const { loading } = this.state

        return (
            <Modal
                title={<I18n id={MessageKey.FrontendUserMenuAccountSettingsTitle} />}
                onCancel={onCancel}
                footer={null}
                open={open}
                destroyOnHidden
            >
                <Preloader loading={loading}>
                    <Tabs
                        className="user-security-settings-tabs"
                        key={initialTab}
                        defaultActiveKey={initialTab}
                        items={[
                            {
                                key: "profile",
                                label: <I18n id={MessageKey.FrontendUserMenuProfileTab} />,
                                children: this.renderProfileTab(),
                                icon: <UserOutlined />,
                            },
                            {
                                key: "password",
                                label: <I18n id={MessageKey.CommonPassword} />,
                                children: this.renderPasswordTab(),
                                icon: <LockOutlined />,
                            },
                            {
                                key: "totp",
                                label: <I18n id={MessageKey.CommonTwoFactorAuthentication} />,
                                children: this.renderTotpTab(),
                                icon: <SafetyOutlined />,
                            },
                        ]}
                    />
                </Preloader>
            </Modal>
        )
    }
}
