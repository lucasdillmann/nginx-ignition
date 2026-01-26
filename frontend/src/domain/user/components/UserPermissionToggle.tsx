import { UserAccessLevel } from "../model/UserAccessLevel"
import { CheckCircleOutlined, EyeOutlined, StopOutlined } from "@ant-design/icons"
import { Flex, Form, Segmented } from "antd"
import type { SegmentedProps } from "antd"
import React from "react"
import "./UserPermissionToggle.css"
import { I18n } from "../../../core/i18n/I18n"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"

export interface UserPermissionToggleProps {
    id: string
    label: MessageKey
    disableReadWrite?: boolean
    disableNoAccess?: boolean
}

export class UserPermissionToggle extends React.Component<UserPermissionToggleProps> {
    private buildOptions(): SegmentedProps["options"] {
        const { disableNoAccess, disableReadWrite } = this.props
        let options = [
            {
                label: <I18n id={MessageKey.FrontendUserComponentsPermissiontoggleNoAccess} />,
                value: UserAccessLevel.NO_ACCESS,
                icon: <StopOutlined />,
            },
            {
                label: <I18n id={MessageKey.FrontendUserComponentsPermissiontoggleReadOnly} />,
                value: UserAccessLevel.READ_ONLY,
                icon: <EyeOutlined />,
            },
            {
                label: <I18n id={MessageKey.FrontendUserComponentsPermissiontoggleFullAccess} />,
                value: UserAccessLevel.READ_WRITE,
                icon: <CheckCircleOutlined />,
            },
        ]

        if (disableReadWrite) {
            options = options.filter(option => option.value !== UserAccessLevel.READ_WRITE)
        }

        if (disableNoAccess) {
            options = options.filter(option => option.value !== UserAccessLevel.NO_ACCESS)
        }

        return options
    }

    render() {
        const { id, label } = this.props

        return (
            <Flex className="user-permission-toggle-container">
                <div className="user-permission-toggle-label">
                    <I18n id={label} />
                </div>
                <div>
                    <Form.Item name={["permissions", id]} style={{ margin: 0, padding: 0 }}>
                        <Segmented options={this.buildOptions()} />
                    </Form.Item>
                </div>
            </Flex>
        )
    }
}
