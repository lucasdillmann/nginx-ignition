import React from "react"
import { Button, Flex, Typography } from "antd"
import { SafetyOutlined } from "@ant-design/icons"
import UserService from "../../../../UserService"
import Notification from "../../../../../../core/components/notification/Notification"
import Preloader from "../../../../../../core/components/preloader/Preloader"
import MessageKey from "../../../../../../core/i18n/model/MessageKey.generated"
import { I18n } from "../../../../../../core/i18n/I18n"
import TotpSetup from "../../../totp/TotpSetup"
import UserConfirmation from "../../../../../../core/components/confirmation/UserConfirmation"
import "./UserTotpTab.css"

interface UserTotpTabState {
    loading: boolean
    enabled?: boolean
}

export default class UserTotpTab extends React.Component<Record<string, never>, UserTotpTabState> {
    private readonly service: UserService

    constructor(props: Record<string, never>) {
        super(props)
        this.service = new UserService()
        this.state = { loading: true }
    }

    componentDidMount() {
        this.fetchStatus()
    }

    private fetchStatus() {
        this.setState({ loading: true })
        this.service
            .getTotpStatus()
            .then(enabled => this.setState({ enabled, loading: false }))
            .catch(() => this.setState({ loading: false }))
    }

    private handleDisabled() {
        UserConfirmation.ask(MessageKey.FrontendUserMenuTotpDisableConfirmation)
            .then(() => this.setState({ loading: true }))
            .then(() => this.service.disableTotp())
            .then(() =>
                Notification.success(
                    MessageKey.FrontendUserMenuTotpDisabledTitle,
                    MessageKey.FrontendUserMenuTotpDisabledSuccessDescription,
                ),
            )
            .catch(() => Notification.error(MessageKey.CommonThatDidntWork, MessageKey.CommonTryAgainLater))
            .then(() => this.setState({ loading: false }))
    }

    render() {
        const { loading, enabled } = this.state

        if (loading) {
            return <Preloader loading={true} />
        }

        if (enabled) {
            return (
                <Flex vertical align="center" justify="center" className="totp-enabled-container">
                    <SafetyOutlined className="totp-enabled-icon" />
                    <Typography.Title level={4}>
                        <I18n id={MessageKey.FrontendUserMenuTotpEnabledTitle} />
                    </Typography.Title>
                    <Typography.Text type="secondary" className="totp-enabled-description">
                        <I18n id={MessageKey.FrontendUserMenuTotpEnabledDescription} />
                    </Typography.Text>
                    <Button danger type="primary" onClick={() => this.handleDisabled()} className="totp-disable-button">
                        <I18n id={MessageKey.FrontendUserMenuTotpDisableButton} />
                    </Button>
                </Flex>
            )
        }

        return <TotpSetup onActivation={() => this.setState({ enabled: true })} />
    }
}
