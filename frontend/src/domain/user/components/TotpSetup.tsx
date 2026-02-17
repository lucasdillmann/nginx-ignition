import React, { createRef } from "react"
import { Button, Input, QRCode, Spin, Typography } from "antd"
import { CopyOutlined, SafetyOutlined } from "@ant-design/icons"
import type { OTPRef } from "antd/es/input/OTP"
import Notification from "../../../core/components/notification/Notification"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { I18n, i18n } from "../../../core/i18n/I18n"
import If from "../../../core/components/flowcontrol/If"
import "./TotpSetup.css"
import UserService from "../UserService"

export interface TotpSetupProps {
    onActivation: () => void
}

interface TotpSetupState {
    loading: boolean
    url?: string
    secret?: string
    code: string
    failed: boolean
    secretCopied: boolean
}

export default class TotpSetup extends React.Component<TotpSetupProps, TotpSetupState> {
    private readonly service: UserService
    private readonly otpRef = createRef<OTPRef>()

    constructor(props: TotpSetupProps) {
        super(props)
        this.service = new UserService()
        this.state = {
            loading: true,
            code: "",
            failed: false,
            secretCopied: false,
        }
    }

    componentDidMount() {
        this.enableTotp()
    }

    private enableTotp() {
        this.setState({ loading: true })
        this.service
            .enableTotp()
            .then(response => {
                const secret = this.extractSecret(response.url)
                this.setState({
                    url: response.url,
                    secret: secret,
                    loading: false,
                })
            })
            .catch(() => {
                this.setState({ loading: false })
            })
    }

    private extractSecret(url: string): string | undefined {
        try {
            const parsed = new URL(url)
            return parsed.searchParams.get("secret") ?? undefined
        } catch {
            return undefined
        }
    }

    private handleCodeChange(value: string) {
        this.setState({ code: value, failed: false }, () => {
            if (value.length === 6) {
                this.activateTotp(value)
            }
        })
    }

    private activateTotp(code: string) {
        this.setState({ loading: true })
        this.service
            .activateTotp(code)
            .then(() => {
                Notification.success(
                    MessageKey.FrontendUserTotpActivatedTitle,
                    MessageKey.FrontendUserTotpActivatedSubtitle,
                )
                this.props.onActivation()
            })
            .catch(() => {
                this.handleActivationFailure()
            })
    }

    private handleActivationFailure() {
        this.setState({ loading: false, code: "", failed: true }, () => {
            this.otpRef.current?.focus()
        })
        Notification.error(
            MessageKey.FrontendAuthenticationTotpFailedTitle,
            MessageKey.FrontendAuthenticationTotpFailedMessage,
        )
    }

    private async handleCopySecret() {
        const { secret } = this.state
        if (!secret) return

        try {
            await navigator.clipboard.writeText(secret)
            this.setState({ secretCopied: true })
            setTimeout(() => this.setState({ secretCopied: false }), 2000)
        } catch {
            // NO-OP
        }
    }

    render() {
        const { loading, url, secret, code, failed } = this.state

        if (loading && !url) {
            return (
                <div className="totp-setup-container">
                    <Spin size="large" />
                </div>
            )
        }

        return (
            <div className="totp-setup-container">
                <div className="totp-setup-header">
                    <SafetyOutlined className="totp-setup-icon" />
                    <Typography.Title level={4} className="totp-setup-title">
                        <I18n id={MessageKey.CommonTwoFactorAuthentication} />
                    </Typography.Title>
                </div>

                <Typography.Text className="totp-setup-subtitle">
                    <I18n id={MessageKey.FrontendUserTotpSubtitle} />
                </Typography.Text>

                <If condition={!!url}>
                    <div className="totp-setup-qr-container">
                        <QRCode value={url!} size={200} color="#000000" />
                    </div>
                </If>

                <If condition={!!secret}>
                    <Button
                        className="totp-setup-secret"
                        onClick={this.handleCopySecret.bind(this)}
                        title={i18n(MessageKey.CommonCopy)}
                        style={{ height: "auto", border: "none" }}
                    >
                        <code className="totp-setup-secret-value">{secret}</code>
                        <CopyOutlined className="totp-setup-secret-copy-icon" />
                    </Button>
                </If>

                <Typography.Text className="totp-setup-verify-prompt">
                    <I18n id={MessageKey.FrontendAuthenticationTotpPrompt} />
                </Typography.Text>

                <Spin spinning={loading}>
                    <Input.OTP
                        ref={this.otpRef}
                        length={6}
                        size="large"
                        value={code}
                        status={failed ? "error" : undefined}
                        formatter={str => str.replace(/\D/g, "")}
                        onChange={this.handleCodeChange.bind(this)}
                    />
                </Spin>
            </div>
        )
    }
}
