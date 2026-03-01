import React, { createRef } from "react"
import { ArrowLeftOutlined, LockOutlined, SafetyOutlined, UserOutlined } from "@ant-design/icons"
import { Navigate } from "react-router-dom"
import { Button, Form, Input, Typography } from "antd"
import type { OTPRef } from "antd/es/input/OTP"
import AppContext from "../../core/components/context/AppContext"
import Preloader from "../../core/components/preloader/Preloader"
import Notification from "../../core/components/notification/Notification"
import "./LoginPage.css"
import UserService from "../user/UserService"
import LoginOutcome from "../user/model/LoginOutcome"
import { navigateTo, queryParams } from "../../core/components/router/AppRouter"
import { LoginFormPage } from "@ant-design/pro-components"
import ThemeContext from "../../core/components/context/ThemeContext"
import LightBackground from "./background/light.jpg"
import DarkBackground from "./background/dark.jpg"
import MessageKey from "../../core/i18n/model/MessageKey.generated"
import { I18n, i18n } from "../../core/i18n/I18n"
import ThemeToggle from "../../core/components/theme/ThemeToggle"
import I18nLanguagePicker from "../../core/i18n/I18nLanguagePicker"
import ShellUserMenuQueue from "../user/components/ShellUserMenuQueue"

interface TotpState {
    failed: boolean
    code: string
}

interface Credentials {
    username: string
    password: string
}

interface LoginPageState {
    loading: boolean
    attemptFailed: boolean
    backgroundImageUrl: string
    totp?: TotpState
    credentials?: Credentials
}

export default class LoginPage extends React.Component<any, LoginPageState> {
    private readonly service: UserService
    private readonly otpRef = createRef<OTPRef>()
    private readonly formRef = createRef<HTMLFormElement>()

    constructor(props: any) {
        super(props)
        this.service = new UserService()
        this.state = {
            loading: false,
            attemptFailed: false,
            backgroundImageUrl: ThemeContext.isDarkMode() ? DarkBackground : LightBackground,
        }
    }

    private handleSubmit(values: { username: string; password: string }) {
        const { username, password } = values
        this.performLogin(username, password)
    }

    private performLogin(username: string, password: string, totp?: string) {
        this.setState({ loading: true })
        this.service
            .login(username, password, totp)
            .then(outcome => this.handleLoginOutcome(outcome, username, password))
            .then(() => this.setState({ loading: false }))
    }

    private async handleLoginOutcome(outcome: LoginOutcome, username: string, password: string) {
        switch (outcome) {
            case LoginOutcome.SUCCESS:
                return this.handleSuccessfulLogin()
            case LoginOutcome.MISSING_TOTP:
                this.setState(
                    {
                        totp: { failed: false, code: "" },
                        credentials: { username, password },
                    },
                    () => {
                        this.otpRef.current?.focus()
                    },
                )
                break
            default:
                if (this.state.totp) {
                    this.setState(
                        {
                            totp: { failed: true, code: "" },
                        },
                        () => this.otpRef.current?.focus(),
                    )
                    Notification.error(
                        MessageKey.FrontendAuthenticationTotpFailedTitle,
                        MessageKey.FrontendAuthenticationTotpFailedMessage,
                    )
                } else {
                    this.handleLoginError()
                }
        }
    }

    private async handleSuccessfulLogin() {
        const { totp } = this.state
        const returnTo = queryParams().returnTo as string | undefined

        return AppContext.get()
            .container!!.reload()
            .then(() => {
                if (returnTo) navigateTo(returnTo, true)
                if (!totp)
                    Notification.warning(
                        MessageKey.FrontendUserMenuTotpDisabledTitle,
                        MessageKey.FrontendUserMenuTotpDisabledWarningDescription,
                        {
                            actions: [
                                <Button
                                    key="show-details"
                                    type="default"
                                    onClick={() => ShellUserMenuQueue.openTotpConfig()}
                                >
                                    <I18n id={MessageKey.FrontendUserMenuTotpDisabledWarningButton} />
                                </Button>,
                            ],
                        },
                    )
            })
    }

    private handleLoginError() {
        this.setState({ attemptFailed: true })
        Notification.error(
            MessageKey.FrontendAuthenticationLoginFailedTitle,
            MessageKey.FrontendAuthenticationLoginFailedMessage,
        )
    }

    private handleTotpSubmit() {
        const { credentials, totp } = this.state
        if (!credentials || totp?.code.length !== 6) return

        this.performLogin(credentials.username, credentials.password, totp.code)
    }

    private handleTotpBack() {
        this.setState(
            {
                totp: undefined,
                credentials: undefined,
            },
            () => this.formRef.current?.resetFields(),
        )
    }

    private renderSubtitle() {
        return (
            <>
                <p style={{ paddingRight: 25 }}>
                    <I18n id={MessageKey.FrontendAuthenticationLoginSubtitle} />
                </p>
                <ThemeToggle />
                <I18nLanguagePicker style={{ marginLeft: 10 }} />
            </>
        )
    }

    private handleTotpChange(value: string) {
        this.setState({ totp: { failed: false, code: value } }, () => {
            if (value.length === 6) {
                this.handleTotpSubmit()
            }
        })
    }

    private renderTotpFields() {
        const { totp } = this.state

        return (
            <div className="totp-container">
                <div className="totp-header">
                    <SafetyOutlined className="totp-icon" />
                    <Typography.Title level={4} className="totp-title">
                        <I18n id={MessageKey.CommonTwoFactorAuthentication} />
                    </Typography.Title>
                </div>
                <Typography.Text className="totp-prompt">
                    <I18n id={MessageKey.FrontendAuthenticationTotpPrompt} />
                </Typography.Text>
                <Input.OTP
                    ref={this.otpRef}
                    length={6}
                    size="large"
                    value={totp?.code}
                    status={totp?.failed ? "error" : undefined}
                    formatter={str => str.replace(/\D/g, "")}
                    onChange={this.handleTotpChange.bind(this)}
                />
                <Button type="link" className="totp-back-link" onClick={this.handleTotpBack.bind(this)}>
                    <ArrowLeftOutlined style={{ marginRight: 6 }} />
                    <I18n id={MessageKey.FrontendAuthenticationTotpBack} />
                </Button>
            </div>
        )
    }

    private renderCredentialFields() {
        return (
            <>
                <Form.Item
                    name="username"
                    validateStatus={this.state.attemptFailed ? "error" : undefined}
                    style={{ marginLeft: 20 }}
                >
                    <Input placeholder={i18n(MessageKey.CommonUsername)} size="large" prefix={<UserOutlined />} />
                </Form.Item>
                <Form.Item
                    name="password"
                    validateStatus={this.state.attemptFailed ? "error" : undefined}
                    style={{ marginLeft: 20 }}
                >
                    <Input.Password
                        placeholder={i18n(MessageKey.CommonPassword)}
                        size="large"
                        prefix={<LockOutlined />}
                    />
                </Form.Item>
            </>
        )
    }

    private renderForm() {
        const { backgroundImageUrl, totp, loading } = this.state

        return (
            <LoginFormPage
                id="nginx-ignition-login-form"
                formRef={this.formRef}
                title={<I18n id={MessageKey.CommonAppName} />}
                subTitle={this.renderSubtitle()}
                onFinish={totp ? this.handleTotpSubmit.bind(this) : this.handleSubmit.bind(this)}
                backgroundImageUrl={backgroundImageUrl}
                submitter={
                    totp
                        ? { render: () => [] }
                        : {
                              searchConfig: {
                                  submitText: <I18n id={MessageKey.FrontendAuthenticationLoginButton} />,
                              },
                              submitButtonProps: {
                                  style: {
                                      width: "auto",
                                      float: "right",
                                  },
                                  loading,
                              },
                          }
                }
                containerStyle={{
                    display: "flex",
                    justifyContent: "center",
                    backgroundColor: "rgba(0, 0, 0, 0.65)",
                    backdropFilter: "blur(4px)",
                    color: "white",
                    padding: "60px 40px 40px",
                }}
                otherStyle={{
                    width: 10,
                }}
            >
                {totp ? this.renderTotpFields() : this.renderCredentialFields()}
            </LoginFormPage>
        )
    }

    private handleThemeChange(darkMode: boolean) {
        const backgroundImageUrl = darkMode ? DarkBackground : LightBackground
        this.setState({ backgroundImageUrl })
    }

    componentDidMount() {
        ThemeContext.register(this.handleThemeChange.bind(this))
    }

    componentWillUnmount() {
        ThemeContext.deregister(this.handleThemeChange.bind(this))
    }

    render() {
        const { loading } = this.state
        const { user, onboardingStatus } = AppContext.get()

        if (user?.id != null) {
            return <Navigate to="/" />
        }

        if (!onboardingStatus.finished) {
            return <Navigate to="/onboarding" />
        }

        return <Preloader loading={loading}>{this.renderForm()}</Preloader>
    }
}
