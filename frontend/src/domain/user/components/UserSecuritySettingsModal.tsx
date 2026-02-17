import React from "react"
import { Button, Flex, Form, FormInstance, Modal, Tabs, Typography } from "antd"
import { LockOutlined, SafetyOutlined } from "@ant-design/icons"
import UserService from "../UserService"
import Notification from "../../../core/components/notification/Notification"
import ValidationResult from "../../../core/validation/ValidationResult"
import Preloader from "../../../core/components/preloader/Preloader"
import UserUpdatePasswordRequest from "../model/UserUpdatePasswordRequest"
import FormLayout from "../../../core/components/form/FormLayout"
import Password from "antd/es/input/Password"
import { UnexpectedResponseError } from "../../../core/apiclient/ApiResponse"
import ValidationResultConverter from "../../../core/validation/ValidationResultConverter"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { I18n } from "../../../core/i18n/I18n"
import TotpSetup from "./TotpSetup"
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import "./UserSecuritySettingsModal.css"

interface UserSecuritySettingsModalProps {
    open: boolean
    onCancel: () => void
    startWithTotp?: boolean
}

interface UserSecuritySettingsModalState {
    loading: boolean
    validationResult: ValidationResult
    passwordFormValues: UserUpdatePasswordRequest
    totpEnabled?: boolean
    totpLoading: boolean
}

const DEFAULT_PASSWORD_FORM_VALUES: UserUpdatePasswordRequest = {
    currentPassword: "",
    newPassword: "",
}

export default class UserSecuritySettingsModal extends React.Component<
    UserSecuritySettingsModalProps,
    UserSecuritySettingsModalState
> {
    private readonly formRef: React.RefObject<FormInstance | null>
    private readonly service: UserService

    constructor(props: UserSecuritySettingsModalProps) {
        super(props)
        this.service = new UserService()
        this.formRef = React.createRef()
        this.state = {
            loading: false,
            validationResult: new ValidationResult(),
            passwordFormValues: DEFAULT_PASSWORD_FORM_VALUES,
            totpLoading: true,
        }
    }

    componentDidUpdate(prevProps: Readonly<UserSecuritySettingsModalProps>) {
        if (this.props.open && !prevProps.open) {
            this.fetchTotpStatus()
        }
    }

    private fetchTotpStatus() {
        this.setState({ totpLoading: true })
        this.service
            .getTotpStatus()
            .then(enabled => this.setState({ totpEnabled: enabled, totpLoading: false }))
            .catch(() => this.setState({ totpLoading: false }))
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
        this.formRef.current?.resetFields()
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

    private renderPasswordTab() {
        const { validationResult, passwordFormValues } = this.state
        return (
            <Form<UserUpdatePasswordRequest>
                {...FormLayout.FormDefaults}
                {...FormLayout.ExpandedLabeledItem}
                ref={this.formRef}
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
        const { open, onCancel, startWithTotp } = this.props
        const { loading } = this.state
        const initialTab = startWithTotp ? "totp" : "password"

        return (
            <Modal
                title={<I18n id={MessageKey.FrontendUserMenuSecuritySettingsTitle} />}
                onCancel={onCancel}
                footer={null}
                open={open}
                destroyOnHidden
            >
                <Preloader loading={loading}>
                    <Tabs
                        defaultActiveKey={initialTab}
                        items={[
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
