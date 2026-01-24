import React from "react"
import { themedModal } from "../theme/ThemedResources"
import { I18n, I18nMessage, i18n } from "../../i18n/I18n"
import MessageKey from "../../i18n/model/MessageKey.generated"

class UserConfirmation {
    ask(message: I18nMessage): Promise<void> {
        return this.askWithCallback(message, () => Promise.resolve())
    }

    askWithCallback<T>(message: I18nMessage, onConfirm: () => Promise<T>): Promise<T> {
        const messageContainer = (
            <div style={{ margin: "0 0 15px" }}>
                <I18n id={message} />
            </div>
        )

        return new Promise(resolve => {
            themedModal().confirm({
                title: i18n(MessageKey.FrontendComponentsConfirmationTitle),
                content: messageContainer,
                cancelText: i18n(MessageKey.CommonNo),
                okText: i18n(MessageKey.CommonYes),
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
