import { themedModal } from "../theme/ThemedResources"

class AccessDeniedModal {
    show() {
        themedModal().error({
            title: "Access denied",
            content: "Sorry, but you can't perform the requested action",
        })
    }
}

export default new AccessDeniedModal()
