import React from "react"
import { themedModal } from "../theme/ThemedResources"

class UserConfirmation {
    ask(message: React.ReactNode): Promise<void> {
        return this.askWithCallback(message, () => Promise.resolve())
    }

    askWithCallback<T>(message: React.ReactNode, onConfirm: () => Promise<T>): Promise<T> {
        const messageContainer = <div style={{ margin: "0 0 15px" }}>{message}</div>

        return new Promise(resolve => {
            themedModal().confirm({
                title: "Beware, young padawan",
                content: messageContainer,
                cancelText: "No",
                okText: "Yes",
                okButtonProps: {
                    color: "danger",
                    variant: "solid",
                },
                onOk: () => onConfirm().then(resolve),
            })
        })
    }
}

export default new UserConfirmation()
