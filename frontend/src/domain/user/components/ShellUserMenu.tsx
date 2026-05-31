import React from "react"
import { Flex } from "antd"
import { LockOutlined, LogoutOutlined, QuestionCircleOutlined } from "@ant-design/icons"
import AppContext from "../../../core/components/context/AppContext"
import ThemeToggle from "../../../core/components/theme/ThemeToggle"
import "./ShellUserMenu.css"
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import UserService from "../UserService"
import { navigateTo } from "../../../core/components/router/AppRouter"
import Notification from "../../../core/components/notification/Notification"
import { buildLoginUrl } from "../../../core/authentication/buildLoginUrl"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import I18nLanguagePicker from "../../../core/i18n/I18nLanguagePicker"
import UserSecuritySettingsModal, { UserSecuritySettingsTab } from "./UserSecuritySettingsModal"
import ShellUserMenuQueue, { QueueAction } from "./ShellUserMenuQueue"

interface ShellUserMenuState {
    modalOpen: boolean
    modalInitialTab: UserSecuritySettingsTab
}

export default class ShellUserMenu extends React.Component<any, ShellUserMenuState> {
    private readonly service: UserService

    constructor(props: any) {
        super(props)
        this.service = new UserService()
        this.state = {
            modalOpen: false,
            modalInitialTab: "password",
        }
    }

    private handleQueueAction(action: QueueAction) {
        if (action == QueueAction.OPEN_TOTP_CONFIG) this.securitySettingsModal(true, "totp")
    }

    private async handleLogout() {
        return UserConfirmation.ask(MessageKey.FrontendUserLogoutConfirmation)
            .then(() => this.service.logout())
            .then(() => Notification.success(MessageKey.CommonSeeYa, MessageKey.FrontendUserLoggedOut))
            .then(() => {
                AppContext.get().user = undefined
            })
            .then(() => navigateTo(buildLoginUrl()))
    }

    private securitySettingsModal(modalOpen: boolean, modalInitialTab: UserSecuritySettingsTab = "password") {
        this.setState({
            modalOpen,
            modalInitialTab,
        })
    }

    componentWillUnmount() {
        ShellUserMenuQueue.detach()
    }

    componentDidMount() {
        ShellUserMenuQueue.attach(action => this.handleQueueAction(action))
    }

    render() {
        const { user } = AppContext.get()
        const { modalOpen, modalInitialTab } = this.state

        return (
            <Flex className="shell-user-menu-container">
                <Flex className="shell-user-menu-actions">
                    <ThemeToggle />
                    <I18nLanguagePicker />
                    <LockOutlined onClick={() => this.securitySettingsModal(true, "password")} />
                    <QuestionCircleOutlined onClick={() => navigateTo("/help")} />
                </Flex>
                <Flex className="shell-user-menu-user-name" onClick={() => this.securitySettingsModal(true, "profile")}>
                    {user?.name}
                </Flex>
                <Flex className="shell-user-menu-icon">
                    <LogoutOutlined onClick={() => this.handleLogout()} />
                </Flex>

                <UserSecuritySettingsModal
                    open={modalOpen}
                    initialTab={modalInitialTab}
                    onCancel={() => this.securitySettingsModal(false)}
                />
            </Flex>
        )
    }
}
