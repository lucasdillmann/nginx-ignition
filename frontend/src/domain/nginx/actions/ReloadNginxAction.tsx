import LocalStorageRepository from "../../../core/repository/LocalStorageRepository"
import { Checkbox, CheckboxRef } from "antd"
import React from "react"
import "./ReloadNginxAction.css"
import GenericNginxAction, { ActionType } from "./GenericNginxAction"
import { themedModal } from "../../../core/components/theme/ThemedResources"
import { isAccessGranted } from "../../../core/components/accesscontrol/IsAccessGranted"
import { UserAccessLevel } from "../../user/model/UserAccessLevel"

class ReloadNginxAction {
    private readonly repository: LocalStorageRepository<boolean>
    private readonly delegate: GenericNginxAction

    constructor() {
        this.repository = new LocalStorageRepository("nginxIgnition.nginxReload.skipConfirmation")
        this.delegate = new GenericNginxAction(ActionType.RELOAD, "nginxIgnition.nginxReloadAction")
    }

    async execute(): Promise<void> {
        if (!isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.nginxServer)) {
            return
        }

        const skipConfirmation = this.repository.getOrDefault(false)
        if (skipConfirmation) return this.delegate.execute()

        return new Promise(resolve => {
            const skipRef = React.createRef<CheckboxRef>()

            themedModal().confirm({
                onCancel: () => resolve(),
                onOk: () => {
                    if (skipRef.current?.input?.checked) this.repository.set(true)

                    return this.delegate.execute().then(resolve)
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
