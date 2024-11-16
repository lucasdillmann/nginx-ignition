import {notification} from "antd";
import {IconType} from "antd/es/notification/interface";

function showNotification(title: string, message: string, type: IconType) {
    notification.open({
        closable: true,
        duration: 5,
        pauseOnHover: true,
        placement: "topRight",
        showProgress: true,
        message: title,
        description: message,
        type,
    })
}

class NotificationFacade {
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

export default new NotificationFacade()
