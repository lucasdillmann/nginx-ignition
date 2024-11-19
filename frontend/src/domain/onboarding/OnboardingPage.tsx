import React from "react";
import {Button, Form, Input, Typography} from "antd";
import {LockOutlined, UserOutlined, IdcardOutlined} from "@ant-design/icons";
import {Navigate} from "react-router-dom";
import AppContext, {AppContextData} from "../../core/components/context/AppContext";
import Preloader from "../../core/components/preloader/Preloader";
import NotificationFacade from "../../core/components/notification/NotificationFacade";
import ValidationResult from "../../core/validation/ValidationResult";
import OnboardingService from "./OnboardingService";
import {UnexpectedResponseError} from "../../core/apiclient/ApiResponse";
import ValidationResultConverter from "../../core/validation/ValidationResultConverter";
import "./OnboardingPage.css"
const {Text, Title} = Typography;

interface OnboardingPageState {
    loading: boolean,
    validationResult: ValidationResult,
    values: any,
}

export default class OnboardingPage extends React.Component<any, OnboardingPageState> {
    static contextType = AppContext
    context!: React.ContextType<typeof AppContext>

    private service: OnboardingService

    constructor(props: any, context: AppContextData) {
        super(props, context);

        this.service = new OnboardingService();
        this.state = {
            loading: false,
            validationResult: new ValidationResult(),
            values: {},
        }
    }

    private handleSubmit(values: { name: string, username: string, password: string }) {
        const {name, username, password} = values

        this.setState({loading: true, values})
        this.service
            .finish(name, username, password)
            .then(() => this.handleSuccess())
            .catch((error) => this.handleError(error))
            .then(() => this.setState({loading: false}))
    }

    private handleSuccess() {
        // eslint-disable-next-line no-restricted-globals
        location.href = "/"
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null)
                this.setState({ validationResult })
        }

        NotificationFacade.error(
            "That didn't work",
            "Please check the form to see if everything seems correct",
        )
    }

    private renderForm() {
        const {validationResult, values} = this.state

        return (
            <section className="onboarding-section">
                <div className="onboarding-container">
                    <div className="onboarding-header">
                    <Title className="onboarding-title">nginx ignition</Title>
                        <Text className="onboarding-text">
                            Welcome to the nginx ignition. This seems to be your first access, and we need to create
                            the first user of the application. Please fill the form below with your details and you
                            will be ready to go.
                        </Text>
                    </div>
                    <Form
                        onFinish={(values) => this.handleSubmit(values)}
                        layout="vertical"
                        requiredMark="optional">
                        <Form.Item
                            name="name"
                            validateStatus={validationResult.getStatus("name")}
                            help={validationResult.getMessage("name")}
                            className="onboarding-form-input"
                            initialValue={values.name}>
                            <Input
                                prefix={<IdcardOutlined />}
                                placeholder="Name"
                                autoFocus
                            />
                        </Form.Item>
                        <Form.Item
                            name="username"
                            validateStatus={validationResult.getStatus("username")}
                            help={validationResult.getMessage("username")}
                            className="onboarding-form-input"
                            initialValue={values.username}>
                            <Input
                                prefix={<UserOutlined />}
                                placeholder="Username"
                            />
                        </Form.Item>
                        <Form.Item
                            name="password"
                            validateStatus={validationResult.getStatus("password")}
                            help={validationResult.getMessage("password")}
                            initialValue={values.password}>
                            <Input.Password
                                prefix={<LockOutlined />}
                                type="password"
                                placeholder="Password"
                            />
                        </Form.Item>
                        <Form.Item>
                            <Button type="primary" htmlType="submit">
                                Continue
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

        if (user?.id != null || onboardingStatus?.finished) {
            return <Navigate to="/" />
        }

        return (
            <Preloader loading={loading}>
                {this.renderForm()}
            </Preloader>
        )
    }
}
