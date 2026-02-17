import React from "react"
import { Button, Steps } from "antd"
import { IdcardOutlined, LockOutlined, SafetyOutlined, UserOutlined } from "@ant-design/icons"
import { Navigate } from "react-router-dom"
import { LoginFormPage, ProFormText } from "@ant-design/pro-components"
import AppContext from "../../core/components/context/AppContext"
import Preloader from "../../core/components/preloader/Preloader"
import Notification from "../../core/components/notification/Notification"
import ValidationResult from "../../core/validation/ValidationResult"
import OnboardingService from "./OnboardingService"
import { UnexpectedResponseError } from "../../core/apiclient/ApiResponse"
import ValidationResultConverter from "../../core/validation/ValidationResultConverter"
import ThemeContext from "../../core/components/context/ThemeContext"
import ThemeToggle from "../../core/components/theme/ThemeToggle"
import I18nLanguagePicker from "../../core/i18n/I18nLanguagePicker"
import TotpSetup from "../user/components/TotpSetup"
import LightBackground from "../authentication/background/light.jpg"
import DarkBackground from "../authentication/background/dark.jpg"
import MessageKey from "../../core/i18n/model/MessageKey.generated"
import { I18n, i18n } from "../../core/i18n/I18n"
import "./OnboardingPage.css"

interface OnboardingPageState {
    loading: boolean
    currentStep: number
    validationResult: ValidationResult
    values: any
    backgroundImageUrl: string
}

export default class OnboardingPage extends React.Component<any, OnboardingPageState> {
    private readonly service: OnboardingService

    constructor(props: any) {
        super(props)

        this.service = new OnboardingService()
        this.state = {
            loading: false,
            currentStep: 0,
            validationResult: new ValidationResult(),
            values: {},
            backgroundImageUrl: ThemeContext.isDarkMode() ? DarkBackground : LightBackground,
        }
    }

    componentDidMount() {
        ThemeContext.register(this.handleThemeChange.bind(this))
    }

    componentWillUnmount() {
        ThemeContext.deregister(this.handleThemeChange.bind(this))
    }

    private handleThemeChange(darkMode: boolean) {
        this.setState({
            backgroundImageUrl: darkMode ? DarkBackground : LightBackground,
        })
    }

    private handleUserFormSubmit(values: { name: string; username: string; password: string }) {
        const { name, username, password } = values
        this.setState({ loading: true, validationResult: new ValidationResult(), values })

        this.service
            .finish(name, username, password)
            .then(() => this.setState({ currentStep: 1, loading: false }))
            .catch(error => this.handleUserFormError(error))
            .then(() => this.setState({ loading: false }))
    }

    private handleUserFormError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const validationResult = ValidationResultConverter.parse(error.response)
            if (validationResult != null) this.setState({ validationResult })
        }

        Notification.error(MessageKey.CommonThatDidntWork, MessageKey.CommonFormCheckMessage)
    }

    private handleTotpCompleted() {
        AppContext.get().container!!.reload()
    }

    private renderSubtitle() {
        return (
            <>
                <p className="onboarding-text-description">
                    <I18n id={MessageKey.FrontendOnboardingSubtitleIntro} />
                    <br />
                    <br />
                    <I18n id={MessageKey.FrontendOnboardingSubtitleInstruction} />
                </p>
                <ThemeToggle />
                <I18nLanguagePicker style={{ marginLeft: 10 }} />
            </>
        )
    }

    private renderSteps() {
        const { currentStep } = this.state

        return (
            <div className="onboarding-steps">
                <Steps
                    current={currentStep}
                    size="small"
                    items={[
                        {
                            title: <I18n id={MessageKey.CommonUser} />,
                            icon: <UserOutlined />,
                        },
                        {
                            title: <I18n id={MessageKey.FrontendOnboardingStepSecurity} />,
                            icon: <SafetyOutlined />,
                        },
                    ]}
                />
            </div>
        )
    }

    private renderUserForm() {
        const { validationResult, values } = this.state

        return (
            <>
                <ProFormText
                    name="name"
                    placeholder={i18n(MessageKey.CommonName)}
                    initialValue={values.name}
                    fieldProps={{
                        size: "large",
                        prefix: <IdcardOutlined />,
                        status: validationResult.getStatus("name") as any,
                    }}
                    help={validationResult.getMessage("name")}
                    style={{ marginLeft: 20 }}
                />
                <ProFormText
                    name="username"
                    placeholder={i18n(MessageKey.CommonUsername)}
                    initialValue={values.username}
                    fieldProps={{
                        size: "large",
                        prefix: <UserOutlined />,
                        status: validationResult.getStatus("username") as any,
                    }}
                    help={validationResult.getMessage("username")}
                    style={{ marginLeft: 20 }}
                />
                <ProFormText.Password
                    name="password"
                    placeholder={i18n(MessageKey.CommonPassword)}
                    initialValue={values.password}
                    fieldProps={{
                        size: "large",
                        prefix: <LockOutlined />,
                        status: validationResult.getStatus("password") as any,
                    }}
                    help={validationResult.getMessage("password")}
                    style={{ marginLeft: 20 }}
                />
            </>
        )
    }

    private renderTotpStep() {
        return (
            <div className="onboarding-totp-step">
                <TotpSetup onActivation={this.handleTotpCompleted.bind(this)} />
                <Button type="link" className="onboarding-skip-button" onClick={this.handleTotpCompleted.bind(this)}>
                    <I18n id={MessageKey.FrontendOnboardingSkipTotp} />
                </Button>
            </div>
        )
    }

    private renderContent() {
        const { currentStep, backgroundImageUrl, loading } = this.state

        if (currentStep === 1) {
            return (
                <LoginFormPage
                    id="nginx-ignition-onboarding-form"
                    title={<I18n id={MessageKey.CommonAppName} />}
                    subTitle={this.renderSubtitle()}
                    backgroundImageUrl={backgroundImageUrl}
                    submitter={{ render: () => [] }}
                    containerStyle={{
                        display: "flex",
                        justifyContent: "center",
                        backgroundColor: "rgba(0, 0, 0, 0.65)",
                        padding: "60px 40px 40px",
                    }}
                    otherStyle={{ width: 10 }}
                >
                    {this.renderSteps()}
                    {this.renderTotpStep()}
                </LoginFormPage>
            )
        }

        return (
            <LoginFormPage
                id="nginx-ignition-onboarding-form"
                title={<I18n id={MessageKey.CommonAppName} />}
                subTitle={this.renderSubtitle()}
                onFinish={this.handleUserFormSubmit.bind(this)}
                backgroundImageUrl={backgroundImageUrl}
                submitter={{
                    searchConfig: {
                        submitText: <I18n id={MessageKey.CommonContinue} />,
                    },
                    submitButtonProps: {
                        style: {
                            width: "auto",
                            float: "right",
                        },
                        loading,
                    },
                }}
                containerStyle={{
                    display: "flex",
                    justifyContent: "center",
                    backgroundColor: "rgba(0, 0, 0, 0.65)",
                    backdropFilter: "blur(4px)",
                    color: "white",
                    padding: "60px 40px 40px",
                }}
                otherStyle={{ width: 10 }}
            >
                {this.renderSteps()}
                {this.renderUserForm()}
            </LoginFormPage>
        )
    }

    render() {
        const { loading } = this.state
        const { user, onboardingStatus } = AppContext.get()

        if (user?.id != null || onboardingStatus?.finished) {
            return <Navigate to="/" />
        }

        return <Preloader loading={loading}>{this.renderContent()}</Preloader>
    }
}
