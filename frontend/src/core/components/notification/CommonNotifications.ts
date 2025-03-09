import Notification from "./Notification"

class CommonNotifications {
    failedToFetch() {
        Notification.error(
            "Unable to fetch the data",
            "We're unable to fetch the data at this time. Please try again later.",
        )
    }
}

export default new CommonNotifications()
