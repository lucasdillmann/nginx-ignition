import LocalStorageRepository from "../../repository/LocalStorageRepository";
import NginxService from "../../../domain/nginx/NginxService";
import Notification from "../notification/Notification";
import {Checkbox, CheckboxRef, Modal} from "antd";
import React from "react";
import "./NginxReload.css"

class NginxReload {
    private readonly repository: LocalStorageRepository<boolean>
    private readonly service: NginxService

    constructor() {
        this.repository = new LocalStorageRepository("nginxIgnition.nginxReload.skipConfirmation")
        this.service = new NginxService()
    }

    private async reload(): Promise<void> {
        const messageTitle = "Reload nginx configuration";
        return this.service
            .reloadConfiguration()
            .then(() => Notification.success(
                messageTitle,
                "Nginx server configuration was reloaded successfully",
            ))
            .catch(() => Notification.error(
                messageTitle,
                "Nginx server failed to reload the configuration. Please check the logs for more details.",
            ))
    }

    async ask(): Promise<void> {
        const skipConfirmation = this.repository.getOrDefault(false)
        if (skipConfirmation)
            return this.reload()

        return new Promise(resolve => {
            const skipRef = React.createRef<CheckboxRef>()

            Modal.confirm({
                onCancel: () => resolve(),
                onOk: () => {
                    if (skipRef.current?.input?.checked)
                        this.repository.set(true)

                    return this.reload().then(resolve)
                },
                type: "confirm",
                title: "Reload nginx configuration?",
                content: (
                    <div className="reload-confirmation-content">
                        <p>In order to apply your changes, we need to reload the nginx server configuration.</p>
                        <Checkbox ref={skipRef}>Always reload automatically and don't ask again</Checkbox>
                    </div>
                ),
                okText: "Reload now",
                cancelText: "Don't reload",
            })
        })
    }
}

// eslint-disable-next-line import/no-anonymous-default-export
export default new NginxReload()