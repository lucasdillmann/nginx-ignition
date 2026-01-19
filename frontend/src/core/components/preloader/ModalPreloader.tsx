import { LoadingOutlined } from "@ant-design/icons"
import React from "react"
import "./ModalPreloader.css"
import { themedModal } from "../theme/ThemedResources"
import { I18n, I18nMessage } from "../../i18n/I18n"

export default class ModalPreloader {
    private timeoutId?: number
    private instance?: { destroy: () => void }

    private open(title: I18nMessage, content: I18nMessage) {
        this.instance = themedModal().info({
            title: <I18n id={title} />,
            content: <I18n id={content} />,
            footer: null,
            icon: <LoadingOutlined className="modal-preloader-spinner" spin />,
        })
        this.timeoutId = undefined
    }

    show(title: I18nMessage, content: I18nMessage) {
        if (this.timeoutId !== undefined) window.clearTimeout(this.timeoutId)

        this.timeoutId = window.setTimeout(() => this.open(title, content), 500)
    }

    close() {
        if (this.timeoutId !== undefined) window.clearTimeout(this.timeoutId)
        else this.instance?.destroy()
    }
}
