import React from "react"
import { LockOutlined, UserOutlined } from "@ant-design/icons"
import { Navigate } from "react-router-dom"
import AppContext from "../../core/components/context/AppContext"
import Preloader from "../../core/components/preloader/Preloader"
import Notification from "../../core/components/notification/Notification"
import "./LoginPage.css"
import UserService from "../user/UserService"
import { navigateTo, queryParams } from "../../core/components/router/AppRouter"
import { LoginFormPage, ProFormText } from "@ant-design/pro-components"
import ThemeContext from "../../core/components/context/ThemeContext"
import LightBackground from "./background/light.jpg"
import DarkBackground from "./background/dark.jpg"
import MessageKey from "../../core/i18n/model/MessageKey.generated"
import { I18n, i18n } from "../../core/i18n/I18n"
import ThemeToggle from "../../core/components/theme/ThemeToggle"
import I18nLanguagePicker from "../../core/i18n/I18nLanguagePicker"

interface LoginPageState {
    loading: boolean
    attemptFailed: boolean
    backgroundImageUrl: string
}

export default class LoginPage extends React.Component<any, LoginPageState> {
    private readonly service: UserService

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

        this.setState({ loading: true })
        this.service
            .login(username, password)
            .then(() => this.handleSuccessfulLogin())
            .catch(() => this.handleLoginError())
            .then(() => this.setState({ loading: false }))
    }

    private async handleSuccessfulLogin() {
        const returnTo = queryParams().returnTo as string | undefined
        return AppContext.get()
            .container!!.reload()
            .then(() => {
                if (returnTo) navigateTo(returnTo, true)
            })
    }

    private handleLoginError() {
        this.setState({ attemptFailed: true })
        Notification.error(
            MessageKey.FrontendAuthenticationLoginFailedTitle,
            MessageKey.FrontendAuthenticationLoginFailedMessage,
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

    private renderForm() {
        const { backgroundImageUrl } = this.state

        return (
            <LoginFormPage
                id="nginx-ignition-login-form"
                title={<I18n id={MessageKey.CommonAppName} />}
                subTitle={this.renderSubtitle()}
                onFinish={this.handleSubmit.bind(this)}
                backgroundImageUrl={backgroundImageUrl}
                submitter={{
                    searchConfig: {
                        submitText: i18n(MessageKey.FrontendAuthenticationLoginButton),
                    },
                }}
                containerStyle={{
                    display: "flex",
                    justifyContent: "center",
                    backgroundColor: "rgba(0, 0, 0, 0.65)",
                    backdropFilter: "blur(4px)",
                    color: "white",
                    padding: "60px 40px",
                }}
                otherStyle={{
                    width: 10,
                }}
            >
                <ProFormText
                    name="username"
                    placeholder={i18n(MessageKey.FrontendAuthenticationUsernamePlaceholder)}
                    fieldProps={{
                        size: "large",
                        prefix: <UserOutlined />,
                    }}
                    style={{
                        marginLeft: 20,
                    }}
                />
                <ProFormText.Password
                    name="password"
                    placeholder={i18n(MessageKey.FrontendAuthenticationPasswordPlaceholder)}
                    fieldProps={{
                        size: "large",
                        prefix: <LockOutlined />,
                    }}
                    style={{
                        marginLeft: 20,
                    }}
                />
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
