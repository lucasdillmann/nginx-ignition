import { App } from "antd"
import type { ModalStaticFunctions } from "antd/es/modal/confirm"
import type { NotificationInstance } from "antd/es/notification/interface"
import ThemeContext from "../context/ThemeContext"

const internalState = {
    notification: {} as NotificationInstance,
    modal: {} as Omit<ModalStaticFunctions, "warn">,
}

export interface ThemedColors {
    DANGER: string
    SUCCESS: string
}

export default () => {
    const staticFunction = App.useApp()

    internalState.modal = staticFunction.modal
    internalState.notification = staticFunction.notification

    return null
}

export function themedModal(): Omit<ModalStaticFunctions, "warn"> {
    return internalState.modal
}

export function themedNotification(): NotificationInstance {
    return internalState.notification
}

export function themedColors(): ThemedColors {
    const darkMode = ThemeContext.isDarkMode()

    return {
        DANGER: darkMode ? "#a81f39" : "#8C2B36",
        SUCCESS: darkMode ? "#6a8737" : "#3d9f07",
    }
}
