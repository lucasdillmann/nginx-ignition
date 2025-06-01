import { LockOutlined } from "@ant-design/icons"
import { Empty } from "antd"
import React from "react"

export default class AccessDeniedPage extends React.Component {
    render() {
        return (
            <Empty
                style={{ marginTop: 50 }}
                image={<LockOutlined style={{ fontSize: 70, color: "var(--nginxIgnition-colorTextDisabled)" }} />}
                description="Sorry, but you don't have access to this page"
            />
        )
    }
}
