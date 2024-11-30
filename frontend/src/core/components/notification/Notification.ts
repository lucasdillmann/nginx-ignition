import { notification } from "antd"
import { IconType } from "antd/es/notification/interface"

function showNotification(title: string, message: string, type: IconType) {
    notification.open({
        closable: true,
        duration: 5,
        pauseOnHover: true,
        placement: "bottomRight",
        showProgress: true,
        message: title,
        description: message,
        type,
    })
}

class Notification {
    warning(title: string, message: string) {
        showNotification(title, message, "warning")
    }

    error(title: string, message: string) {
        showNotification(title, message, "error")
    }

    success(title: string, message: string) {
        showNotification(title, message, "success")
    }
}

// eslint-disable-next-line import/no-anonymous-default-export
export default new Notification()
