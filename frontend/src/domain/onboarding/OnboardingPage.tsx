import React from "react"
import { Button, Form, Input, Typography } from "antd"
import { IdcardOutlined, LockOutlined, UserOutlined } from "@ant-design/icons"
import { Navigate } from "react-router-dom"
import AppContext from "../../core/components/context/AppContext"
import Preloader from "../../core/components/preloader/Preloader"
import Notification from "../../core/components/notification/Notification"
import ValidationResult from "../../core/validation/ValidationResult"
import OnboardingService from "./OnboardingService"
import { UnexpectedResponseError } from "../../core/apiclient/ApiResponse"
import ValidationResultConverter from "../../core/validation/ValidationResultConverter"
import "./OnboardingPage.css"
import MessageKey from "../../core/i18n/model/MessageKey.generated"

const { Text, Title } = Typography

interface OnboardingPageState {
    loading: boolean
    validationResult: ValidationResult
    values: any
}

export default class OnboardingPage extends React.Component<any, OnboardingPageState> {
    private readonly service: OnboardingService

    constructor(props: any) {
        super(props)

        this.service = new OnboardingService()
        this.state = {
            loading: false,
            validationResult: new ValidationResult(),
            values: {},
        }
    }

    private handleSubmit(values: { name: string; username: string; password: string }) {
        const { name, username, password } = values

        this.setState({ loading: true, validationResult: new ValidationResult(), values })
        this.service
            .finish(name, username, password)
            .then(() => this.handleSuccess())
            .catch(error => this.handleError(error))
            .then(() => this.setState({ loading: false }))
    }

    private async handleSuccess() {
        return AppContext.get().container!!.reload()
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null) this.setState({ validationResult })
        }

        Notification.error(MessageKey.CommonThatDidntWork, MessageKey.CommonFormCheckMessage)
    }

    private renderForm() {
        const { validationResult, values } = this.state

        return (
            <section className="onboarding-section">
                <div className="onboarding-container">
                    <div className="onboarding-header">
                        <Title className="onboarding-title">nginx ignition</Title>
                        <Text className="onboarding-text">
                            Welcome to the nginx ignition. This seems to be your first access, please fill the form
                            below to the create your user and we'll be ready to go.
                        </Text>
                    </div>
                    <Form onFinish={values => this.handleSubmit(values)} layout="vertical" requiredMark="optional">
                        <Form.Item
                            name="name"
                            validateStatus={validationResult.getStatus("name")}
                            help={validationResult.getMessage("name")}
                            className="onboarding-form-input"
                            initialValue={values.name}
                        >
                            <Input size="large" prefix={<IdcardOutlined />} placeholder="Name" autoFocus />
                        </Form.Item>
                        <Form.Item
                            name="username"
                            validateStatus={validationResult.getStatus("username")}
                            help={validationResult.getMessage("username")}
                            className="onboarding-form-input"
                            initialValue={values.username}
                        >
                            <Input size="large" prefix={<UserOutlined />} placeholder="Username" />
                        </Form.Item>
                        <Form.Item
                            name="password"
                            validateStatus={validationResult.getStatus("password")}
                            help={validationResult.getMessage("password")}
                            initialValue={values.password}
                        >
                            <Input.Password
                                size="large"
                                prefix={<LockOutlined />}
                                type="password"
                                placeholder="Password"
                            />
                        </Form.Item>
                        <Form.Item>
                            <Button size="large" type="primary" htmlType="submit">
                                Continue
                            </Button>
                        </Form.Item>
                    </Form>
                </div>
            </section>
        )
    }

    render() {
        const { loading } = this.state
        const { user, onboardingStatus } = AppContext.get()

        if (user?.id != null || onboardingStatus?.finished) {
            return <Navigate to="/" />
        }

        return <Preloader loading={loading}>{this.renderForm()}</Preloader>
    }
}
