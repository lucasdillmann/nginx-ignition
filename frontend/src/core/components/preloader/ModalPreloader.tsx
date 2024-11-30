import { Modal } from "antd"
import { LoadingOutlined } from "@ant-design/icons"
import React from "react"
import "./ModalPreloader.css"

export default class ModalPreloader {
    private timeoutId?: number
    private instance?: { destroy: () => void }

    private open(title: string, content: string) {
        this.instance = Modal.info({
            title,
            content,
            footer: null,
            icon: <LoadingOutlined className="modal-preloader-spinner" spin />,
        })
        this.timeoutId = undefined
    }

    show(title: string, content: string) {
        if (this.timeoutId !== undefined) window.clearTimeout(this.timeoutId)

        this.timeoutId = window.setTimeout(() => this.open(title, content), 500)
    }

    close() {
        if (this.timeoutId !== undefined) window.clearTimeout(this.timeoutId)
        else this.instance?.destroy()
    }
}
