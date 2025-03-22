import LocalStorageRepository from "../../../core/repository/LocalStorageRepository"
import NginxService from "../NginxService"
import Notification from "../../../core/components/notification/Notification"
import { Checkbox, CheckboxRef, Modal } from "antd"
import React from "react"
import "./ReloadNginxAction.css"

const MESSAGE_KEY = "nginxIgnition.nginxReload.message"
const MESSAGE_TITLE = "Reload nginx configuration"

class ReloadNginxAction {
    private readonly repository: LocalStorageRepository<boolean>
    private readonly service: NginxService
    private running: boolean

    constructor() {
        this.repository = new LocalStorageRepository("nginxIgnition.nginxReload.skipConfirmation")
        this.service = new NginxService()
        this.running = false
    }

    private async reload(): Promise<void> {
        if (this.running) return Promise.resolve()

        this.running = true
        this.showInProgressNotification()

        return this.service
            .reloadConfiguration()
            .then(() => this.showSuccessNotification())
            .catch(() => this.showErrorNotification())
            .then(() => {
                this.running = false
            })
    }

    private showInProgressNotification() {
        Notification.progress(MESSAGE_TITLE, "Please wait while we reload the nginx configuration...", MESSAGE_KEY)
    }

    private showSuccessNotification() {
        Notification.success(MESSAGE_TITLE, "Nginx server configuration was reloaded successfully", MESSAGE_KEY)
    }

    private showErrorNotification() {
        Notification.error(
            MESSAGE_TITLE,
            "Nginx server failed to reload the configuration. Please check the logs for more details.",
            MESSAGE_KEY,
        )
    }

    async execute(): Promise<void> {
        const skipConfirmation = this.repository.getOrDefault(false)
        if (skipConfirmation) return this.reload()

        return new Promise(resolve => {
            const skipRef = React.createRef<CheckboxRef>()

            Modal.confirm({
                onCancel: () => resolve(),
                onOk: () => {
                    if (skipRef.current?.input?.checked) this.repository.set(true)

                    return this.reload().then(resolve)
                },
                type: "confirm",
                title: "Reload nginx configuration?",
                content: (
                    <div className="reload-confirmation-content">
                        <p>In order to apply your changes, we need to reload the nginx server configuration.</p>
                        <Checkbox ref={skipRef}>Always reload automatically</Checkbox>
                    </div>
                ),
                okText: "Reload now",
                cancelText: "Don't reload",
            })
        })
    }
}

export default new ReloadNginxAction()
