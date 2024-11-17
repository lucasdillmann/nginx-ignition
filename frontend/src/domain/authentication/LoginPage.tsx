import React from "react";
import {Button, Form, Input, Typography} from "antd";
import {LockOutlined, UserOutlined} from "@ant-design/icons";
import LoginService from "./LoginService";
import {Navigate} from "react-router-dom";
import AppContext, {AppContextData} from "../../core/components/context/AppContext";
import Preloader from "../../core/components/preloader/Preloader";
import NotificationFacade from "../../core/components/notification/NotificationFacade";

const {Text, Title} = Typography;
const styles = {
    container: {
        margin: "0 auto",
        width: "380px"
    },
    section: {
        alignItems: "center",
        display: "flex",
    },
    title: {
        marginTop: 0,
        fontSize: 24,
    },
    header: {
        marginTop: 40,
        marginBottom: 30,
    },
    text: {
        fontSize: 14,
    },
    formInput: {
        marginBottom: 10,
    }
}

interface LoginPageState {
    loading: boolean,
    attemptFailed: boolean,
}

export default class LoginPage extends React.Component<any, LoginPageState> {
    static contextType = AppContext
    context!: React.ContextType<typeof AppContext>

    private service: LoginService

    constructor(props: any, context: AppContextData) {
        super(props, context);

        this.service = new LoginService();
        this.state = {
            loading: false,
            attemptFailed: false,
        }
    }

    private handleSubmit(values: { username: string, password: string }) {
        const {username, password} = values

        this.setState({loading: true})
        this.service
            .login(username, password)
            .then(() => this.handleSuccessfulLogin())
            .catch(() => this.handleLoginError())
            .then(() => this.setState({loading: false}))
    }

    private handleSuccessfulLogin() {
        // eslint-disable-next-line no-restricted-globals
        location.href = "/"
    }

    private handleLoginError() {
        this.setState({attemptFailed: true})
        NotificationFacade.error(
            "Login failed",
            "Please check your username and password.",
        )
    }

    private renderForm() {
        const {attemptFailed} = this.state
        const inputStatus = attemptFailed ? "error" : undefined

        return (
            <section style={styles.section}>
                <div style={styles.container}>
                    <div style={styles.header}>
                        <Title style={styles.title}>nginx ignition</Title>
                        <Text style={styles.text}>
                            Welcome back. Please enter your username and password to continue.
                        </Text>
                    </div>
                    <Form
                        onFinish={(values) => this.handleSubmit(values)}
                        layout="vertical"
                        requiredMark="optional"
                    >
                        <Form.Item name="username" style={styles.formInput}>
                            <Input
                                prefix={<UserOutlined/>}
                                placeholder="Username"
                                status={inputStatus}
                                autoFocus
                            />
                        </Form.Item>
                        <Form.Item name="password">
                            <Input.Password
                                prefix={<LockOutlined/>}
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
        const {loading} = this.state
        const {user, onboardingStatus} = this.context

        if (user?.id != null) {
            return <Navigate to="/" />
        }

        if (!onboardingStatus.finished) {
            return <Navigate to="/onboarding" />
        }

        if (loading) {
            return (
                <Preloader>
                    {this.renderForm()}
                </Preloader>
            )
        }

        return this.renderForm()
    }
}
