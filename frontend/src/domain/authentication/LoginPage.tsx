import React from "react"
import { Button, Form, Input, Typography } from "antd"
import { LockOutlined, UserOutlined } from "@ant-design/icons"
import { Navigate } from "react-router-dom"
import AppContext, { AppContextData } from "../../core/components/context/AppContext"
import Preloader from "../../core/components/preloader/Preloader"
import Notification from "../../core/components/notification/Notification"
import "./LoginPage.css"
import UserService from "../user/UserService"

const { Text, Title } = Typography

interface LoginPageState {
    loading: boolean
    attemptFailed: boolean
}

export default class LoginPage extends React.Component<any, LoginPageState> {
    static contextType = AppContext
    context!: React.ContextType<typeof AppContext>

    private service: UserService

    constructor(props: any, context: AppContextData) {
        super(props, context)

        this.service = new UserService()
        this.state = {
            loading: false,
            attemptFailed: false,
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

    private handleSuccessfulLogin() {
        // eslint-disable-next-line no-restricted-globals
        location.href = "/"
    }

    private handleLoginError() {
        this.setState({ attemptFailed: true })
        Notification.error("Login failed", "Please check your username and password.")
    }

    private renderForm() {
        const { attemptFailed } = this.state
        const inputStatus = attemptFailed ? "error" : undefined

        return (
            <section className="login-section">
                <div className="login-container">
                    <div className="login-header">
                        <Title className="login-title">nginx ignition</Title>
                        <Text className="login-text">
                            Welcome back. Please enter your username and password to continue.
                        </Text>
                    </div>
                    <Form onFinish={values => this.handleSubmit(values)} layout="vertical" requiredMark="optional">
                        <Form.Item name="username" className="login-form-input">
                            <Input prefix={<UserOutlined />} placeholder="Username" status={inputStatus} autoFocus />
                        </Form.Item>
                        <Form.Item name="password">
                            <Input.Password
                                prefix={<LockOutlined />}
                                type="password"
                                placeholder="Password"
                                status={inputStatus}
                            />
                        </Form.Item>
                        <Form.Item>
                            <Button type="primary" htmlType="submit">
                                Log in
                            </Button>
                        </Form.Item>
                    </Form>
                </div>
            </section>
        )
    }

    render() {
        const { loading } = this.state
        const { user, onboardingStatus } = this.context

        if (user?.id != null) {
            return <Navigate to="/" />
        }

        if (!onboardingStatus.finished) {
            return <Navigate to="/onboarding" />
        }

        return <Preloader loading={loading}>{this.renderForm()}</Preloader>
    }
}
