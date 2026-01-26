import { ExclamationCircleOutlined } from "@ant-design/icons"
import { Empty } from "antd"
import React from "react"
import { I18n } from "../../i18n/I18n"
import MessageKey from "../../i18n/model/MessageKey.generated"

class EmptyStates {
    public FailedToFetch = (
        <Empty
            image={
                <ExclamationCircleOutlined style={{ fontSize: 70, color: "var(--nginxIgnition-colorTextDisabled)" }} />
            }
            description={<I18n id={MessageKey.FrontendComponentsEmptystateFailedToFetch} />}
        />
    )

    public NotFound = (<Empty description={<I18n id={MessageKey.CommonNotFoundTitle} />} />)
}

export default new EmptyStates()
