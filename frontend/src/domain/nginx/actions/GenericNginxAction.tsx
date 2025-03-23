import NginxService from "../NginxService"
import Notification, { Props } from "../../../core/components/notification/Notification"
import React from "react"
import { Button, Modal } from "antd"
import { NginxActionResponse } from "../model/NginxActionResponse"

export enum ActionType {
    RELOAD,
    START,
    STOP,
}

interface Message {
    title: string
    content: string
}

interface Action {
    action: () => Promise<NginxActionResponse>
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
                            title: "Reload nginx server configuration",
                            content: "Please wait while we reload the nginx configuration",
                        },
                        success: {
                            title: "Reload nginx server configuration",
                            content: "Nginx configuration was reloaded successfully",
                        },
                        error: {
                            title: "Reload nginx server configuration failed",
                            content: "We're unable to reload the nginx configuration at this moment",
                        },
                    },
                }
            case ActionType.START:
                return {
                    action: () => this.service.start(),
                    messages: {
                        inProgress: {
                            title: "Start nginx server",
                            content: "Please wait while we start the nginx server",
                        },
                        success: {
                            title: "Start nginx server",
                            content: "Nginx server was started successfully",
                        },
                        error: {
                            title: "Start nginx server failed",
                            content: "We're unable to start the nginx server at this moment",
                        },
                    },
                }
            case ActionType.STOP:
                return {
                    action: () => this.service.stop(),
                    messages: {
                        inProgress: {
                            title: "Stop nginx server",
                            content: "Please wait while we stop the nginx server",
                        },
                        success: {
                            title: "Stop nginx server",
                            content: "Nginx server was stopped successfully",
                        },
                        error: {
                            title: "Stop nginx server failed",
                            content: "We're unable to stop the nginx server at this moment",
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
            Modal.error({
                width: 750,
                title: "Error details",
                content: (
                    <>
                        <p>The following message was produced by the nginx server with the details of what happened.</p>
                        <code>{error.response?.body?.message ?? error.message}</code>
                    </>
                ),
            })

        Notification.error(title, content, {
            ...this.messageProps(),
            actions: [
                <Button key="show-details" type="default" onClick={onClick}>
                    Open error details
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
        if (this.startedAt !== undefined) return Promise.resolve()

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
