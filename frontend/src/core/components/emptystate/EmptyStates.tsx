import { ExclamationCircleOutlined } from "@ant-design/icons"
import { Empty } from "antd"
import React from "react"

class EmptyStates {
    public FailedToFetch = (
        <Empty
            image={
                <ExclamationCircleOutlined style={{ fontSize: 70, color: "var(--nginxIgnition-colorTextDisabled)" }} />
            }
            description="Unable to fetch the data. Please try again later."
        />
    )

    public NotFound = (<Empty description="Not found" />)
}

export default new EmptyStates()
