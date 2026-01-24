import Notification from "./Notification"
import MessageKey from "../../i18n/model/MessageKey.generated"

class CommonNotifications {
    failedToFetch() {
        Notification.error(MessageKey.CommonUnableToFetchData, MessageKey.CommonTryAgainLater)
    }
}

export default new CommonNotifications()
