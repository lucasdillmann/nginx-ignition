import { UserAccessLevel } from "../model/UserAccessLevel"
import { CheckCircleOutlined, EyeOutlined, StopOutlined } from "@ant-design/icons"
import { Flex, Segmented } from "antd"
import React from "react"
import "./UserPermissionToggle.css"

export interface UserPermissionToggleProps {
    label: string
    value: UserAccessLevel
    onChange: (value: UserAccessLevel) => void
    disableReadWrite?: boolean
    disableNoAccess?: boolean
}

const NO_ACCESS_OPTION = {
    label: "No access",
    value: UserAccessLevel.NO_ACCESS,
    icon: <StopOutlined />,
}

const READ_ONLY_OPTION = {
    label: "Read only",
    value: UserAccessLevel.READ_ONLY,
    icon: <EyeOutlined />,
}

const READ_WRITE_OPTION = {
    label: "Full access",
    value: UserAccessLevel.READ_WRITE,
    icon: <CheckCircleOutlined />,
}

export class UserPermissionToggle extends React.Component<UserPermissionToggleProps> {
    private buildOptions() {
        const { disableNoAccess, disableReadWrite } = this.props
        let options = [NO_ACCESS_OPTION, READ_ONLY_OPTION, READ_WRITE_OPTION]

        if (disableReadWrite) {
            options = options.filter(option => option != READ_WRITE_OPTION)
        }

        if (disableNoAccess) {
            options = options.filter(option => option != NO_ACCESS_OPTION)
        }

        return options
    }

    render() {
        const { onChange, value, label } = this.props

        return (
            <Flex className="user-permission-toggle-container">
                <div className="user-permission-toggle-label">{label}</div>
                <div>
                    <Segmented label={label} options={this.buildOptions()} value={value} onChange={onChange} />
                </div>
            </Flex>
        )
    }
}
