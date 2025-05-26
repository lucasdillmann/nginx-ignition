import { App } from "antd"
import type { ModalStaticFunctions } from "antd/es/modal/confirm"
import type { NotificationInstance } from "antd/es/notification/interface"

const internalState = {
    notification: {} as NotificationInstance,
    modal: {} as Omit<ModalStaticFunctions, "warn">,
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
