import { IconType } from "antd/es/notification/interface"
import React from "react"
import { LoadingOutlined } from "@ant-design/icons"
import { themedNotification } from "../theme/ThemedResources"
import { I18n, I18nMessage } from "../../i18n/I18n"

export interface Props {
    key?: React.Key
    duration?: number
    actions?: React.ReactNode[]
}

function showNotification(
    title: I18nMessage,
    message: I18nMessage,
    type: IconType,
    props?: Props,
    role?: "alert" | "status",
    icon?: React.ReactNode,
) {
    let duration: number | undefined
    if (props?.duration) {
        duration = props.duration
    } else {
        duration = role === "status" ? 0 : 5
    }

    themedNotification().open({
        closable: role !== "status",
        showProgress: role !== "status",
        placement: "bottomRight",
        pauseOnHover: true,
        description: <I18n id={message} />,
        key: props?.key,
        actions: props?.actions,
        message: <I18n id={title} />,
        duration,
        type,
        role,
        icon,
    })
}

class Notification {
    warning(title: I18nMessage, message: I18nMessage, props?: Props) {
        showNotification(title, message, "warning", props)
    }

    error(title: I18nMessage, message: I18nMessage, props?: Props) {
        showNotification(title, message, "error", props)
    }

    success(title: I18nMessage, message: I18nMessage, props?: Props) {
        showNotification(title, message, "success", props)
    }

    progress(title: I18nMessage, message: I18nMessage, props?: Props) {
        const icon = <LoadingOutlined style={{ fontSize: 24 }} spin />
        showNotification(title, message, "info", props, "status", icon)
    }

    close(key: React.Key) {
        themedNotification().destroy(key)
    }
}

export default new Notification()
