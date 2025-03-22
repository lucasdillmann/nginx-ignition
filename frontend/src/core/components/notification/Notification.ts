import { notification } from "antd"
import { IconType } from "antd/es/notification/interface"
import React from "react"

function showNotification(title: string, message: string, type: IconType, key?: React.Key, role?: "alert" | "status") {
    notification.open({
        closable: role !== "status",
        duration: role === "status" ? undefined : 5,
        pauseOnHover: role !== "status",
        placement: "bottomRight",
        showProgress: role !== "status",
        message: title,
        description: message,
        type,
        key,
        role,
    })
}

class Notification {
    warning(title: string, message: string, key?: React.Key) {
        showNotification(title, message, "warning", key)
    }

    error(title: string, message: string, key?: React.Key) {
        showNotification(title, message, "error", key)
    }

    success(title: string, message: string, key?: React.Key) {
        showNotification(title, message, "success", key)
    }

    progress(title: string, message: string, key?: React.Key) {
        showNotification(title, message, "info", key, "status")
    }

    close(key: React.Key) {
        notification.destroy(key)
    }
}

export default new Notification()
