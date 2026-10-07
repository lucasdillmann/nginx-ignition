import React from "react"
import { Modal, Tabs } from "antd"
import { KeyOutlined, LockOutlined, SafetyOutlined, UserOutlined } from "@ant-design/icons"
import MessageKey from "../../../../core/i18n/model/MessageKey.generated"
import { I18n } from "../../../../core/i18n/I18n"
import UserProfileTab from "./tabs/profile/UserProfileTab"
import UserPasswordTab from "./tabs/password/UserPasswordTab"
import UserTotpTab from "./tabs/totp/UserTotpTab"
import UserAPITokensTab from "./tabs/token/UserAPITokensTab"
import "./UserAccountSettingsModal.css"

export type UserAccountSettingsTab = "profile" | "password" | "totp" | "tokens"

interface UserAccountSettingsModalProps {
    open: boolean
    onCancel: () => void
    initialTab?: UserAccountSettingsTab
}

export default class UserAccountSettingsModal extends React.Component<UserAccountSettingsModalProps> {
    render() {
        const { open, onCancel, initialTab = "password" } = this.props

        return (
            <Modal
                title={<I18n id={MessageKey.FrontendUserMenuAccountSettingsTitle} />}
                onCancel={onCancel}
                footer={null}
                open={open}
                width={700}
                destroyOnHidden
            >
                <Tabs
                    className="user-account-settings-tabs"
                    key={initialTab}
                    defaultActiveKey={initialTab}
                    items={[
                        {
                            key: "profile",
                            label: <I18n id={MessageKey.FrontendUserMenuProfileTab} />,
                            children: <UserProfileTab onClose={onCancel} />,
                            icon: <UserOutlined />,
                        },
                        {
                            key: "password",
                            label: <I18n id={MessageKey.CommonPassword} />,
                            children: <UserPasswordTab onClose={onCancel} />,
                            icon: <LockOutlined />,
                        },
                        {
                            key: "totp",
                            label: <I18n id={MessageKey.FrontendUserMenuTotpTabTitle} />,
                            children: <UserTotpTab />,
                            icon: <SafetyOutlined />,
                        },
                        {
                            key: "tokens",
                            label: <I18n id={MessageKey.FrontendUserTokensTabTitle} />,
                            children: <UserAPITokensTab />,
                            icon: <KeyOutlined />,
                        },
                    ]}
                />
            </Modal>
        )
    }
}
