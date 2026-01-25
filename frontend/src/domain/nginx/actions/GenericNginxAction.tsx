import NginxService from "../NginxService"
import Notification, { Props } from "../../../core/components/notification/Notification"
import React from "react"
import { Button } from "antd"
import { themedModal } from "../../../core/components/theme/ThemedResources"
import { UserAccessLevel } from "../../user/model/UserAccessLevel"
import { isAccessGranted } from "../../../core/components/accesscontrol/IsAccessGranted"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { I18n, I18nMessage } from "../../../core/i18n/I18n"

export enum ActionType {
    RELOAD,
    START,
    STOP,
}

interface Message {
    title: I18nMessage
    content: I18nMessage
}

interface Action {
    action: () => Promise<void>
    messages: {
        inProgress: Message
        success: Message
        error: Message
    }
}

export default class GenericNginxAction {
    private readonly service: NginxService
    private readonly type: ActionType
    private readonly key: string
    private startedAt?: Date

    constructor(type: ActionType, key: string) {
        this.service = new NginxService()
        this.type = type
        this.key = key
        this.startedAt = undefined
    }

    private metadata(): Action {
        switch (this.type) {
            case ActionType.RELOAD:
                return {
                    action: () => this.service.reloadConfiguration(),
                    messages: {
                        inProgress: {
                            title: MessageKey.FrontendNginxReloadInProgress,
                            content: MessageKey.FrontendNginxReloadInProgressContent,
                        },
                        success: {
                            title: MessageKey.FrontendNginxReloadSuccess,
                            content: MessageKey.FrontendNginxReloadSuccessDescription,
                        },
                        error: {
                            title: MessageKey.FrontendNginxReloadError,
                            content: MessageKey.FrontendNginxReloadErrorDescription,
                        },
                    },
                }
            case ActionType.START:
                return {
                    action: () => this.service.start(),
                    messages: {
                        inProgress: {
                            title: MessageKey.FrontendNginxStartInProgress,
                            content: MessageKey.FrontendNginxStartInProgressContent,
                        },
                        success: {
                            title: MessageKey.FrontendNginxStartSuccess,
                            content: MessageKey.FrontendNginxStartSuccessDescription,
                        },
                        error: {
                            title: MessageKey.FrontendNginxStartError,
                            content: MessageKey.FrontendNginxStartErrorDescription,
                        },
                    },
                }
            case ActionType.STOP:
                return {
                    action: () => this.service.stop(),
                    messages: {
                        inProgress: {
                            title: MessageKey.FrontendNginxStopInProgress,
                            content: MessageKey.FrontendNginxStopInProgressContent,
                        },
                        success: {
                            title: MessageKey.FrontendNginxStopSuccess,
                            content: MessageKey.FrontendNginxStopSuccessDescription,
                        },
                        error: {
                            title: MessageKey.FrontendNginxStopError,
                            content: MessageKey.FrontendNginxStopErrorDescription,
                        },
                    },
                }
        }
    }

    private showInProgressNotification() {
        const { title, content } = this.metadata().messages.inProgress
        Notification.progress(title, content, { key: this.key })
    }

    private showSuccessNotification() {
        const { title, content } = this.metadata().messages.success
        Notification.success(title, content, this.messageProps())
    }

    private showErrorNotification(error: any) {
        const { title, content } = this.metadata().messages.error
        const onClick = () =>
            themedModal().error({
                width: 750,
                title: <I18n id={MessageKey.FrontendComponentsErrorDetails} />,
                content: <code>{error.response?.body?.message ?? error.message}</code>,
            })

        Notification.error(title, content, {
            ...this.messageProps(),
            actions: [
                <Button key="show-details" type="default" onClick={onClick}>
                    <I18n id={MessageKey.CommonOpenErrorDetails} />
                </Button>,
            ],
        })
    }

    private messageProps(): Props {
        const timeDiff = (new Date().getTime() - this.startedAt!!.getTime()) / 1000
        return {
            key: this.key,
            duration: 5 + timeDiff,
        }
    }

    async execute(): Promise<void> {
        const readOnly = !isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.nginxServer)
        if (readOnly || this.startedAt !== undefined) return Promise.resolve()

        this.startedAt = new Date()
        this.showInProgressNotification()

        return this.metadata()
            .action()
            .then(() => this.showSuccessNotification())
            .catch(error => this.showErrorNotification(error))
            .then(() => {
                this.startedAt = undefined
            })
    }
}
