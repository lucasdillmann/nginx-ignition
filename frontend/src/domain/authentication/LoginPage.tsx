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

interface LoginPageState {
    loading: boolean
    attemptFailed: boolean
}

export default class LoginPage extends React.Component<any, LoginPageState> {
    private readonly service: UserService
    private readonly backgroundImageUrl: string

    constructor(props: any) {
        super(props)
        this.service = new UserService()
        this.state = {
            loading: false,
            attemptFailed: false,
        }

        this.backgroundImageUrl = ThemeContext.isDarkMode() ? DarkBackground : LightBackground
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
        Notification.error("Login failed", "Please check your username and password.")
    }

    private renderForm() {
        return (
            <LoginFormPage
                id="nginx-ignition-login-form"
                title="nginx ignition"
                subTitle="Welcome back. Please sign in to continue."
                onFinish={this.handleSubmit.bind(this)}
                backgroundImageUrl={this.backgroundImageUrl}
                submitter={{
                    searchConfig: {
                        submitText: "Log in",
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
                    placeholder="username"
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
                    placeholder="password"
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
