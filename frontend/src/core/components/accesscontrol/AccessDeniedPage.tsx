import { LockOutlined } from "@ant-design/icons"
import { Empty } from "antd"
import React from "react"
import { I18n } from "../../i18n/I18n"
import MessageKey from "../../i18n/model/MessageKey.generated"

export default class AccessDeniedPage extends React.Component {
    render() {
        return (
            <Empty
                style={{ marginTop: 50 }}
                image={<LockOutlined style={{ fontSize: 70, color: "var(--nginxIgnition-colorTextDisabled)" }} />}
                description={<I18n id={MessageKey.FrontendComponentsAccesscontrolAccessDeniedDescription} />}
            />
        )
    }
}
