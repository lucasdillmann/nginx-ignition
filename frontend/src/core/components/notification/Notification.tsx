import { notification } from "antd"
import { IconType } from "antd/es/notification/interface"
import React from "react"
import { LoadingOutlined } from "@ant-design/icons"

export interface Props {
    key?: React.Key
    duration?: number
    actions?: React.ReactNode[]
}

function showNotification(
    title: string,
    message: string,
    type: IconType,
    props?: Props,
    role?: "alert" | "status",
    icon?: React.ReactNode,
) {
    let duration: number | undefined
    if (props?.duration) {
        duration = props.duration
    } else {
        duration = role === "status" ? undefined : 5
    }

    notification.open({
        closable: role !== "status",
        showProgress: role !== "status",
        placement: "bottomRight",
        pauseOnHover: true,
        message: title,
        description: message,
        key: props?.key,
        actions: props?.actions,
        duration,
        type,
        role,
        icon,
    })
}

class Notification {
    warning(title: string, message: string, props?: Props) {
        showNotification(title, message, "warning", props)
    }

    error(title: string, message: string, props?: Props) {
        showNotification(title, message, "error", props)
    }

    success(title: string, message: string, props?: Props) {
        showNotification(title, message, "success", props)
    }

    progress(title: string, message: string, props?: Props) {
        const icon = <LoadingOutlined style={{ fontSize: 24 }} spin />
        showNotification(title, message, "info", props, "status", icon)
    }

    close(key: React.Key) {
        notification.destroy(key)
    }
}

export default new Notification()
