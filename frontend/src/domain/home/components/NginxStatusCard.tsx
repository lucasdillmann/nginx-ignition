import React from "react"
import { Button, ConfigProvider, Flex } from "antd"
import Preloader from "../../../core/components/preloader/Preloader"
import NginxService from "../../nginx/NginxService"
import NginxEventDispatcher from "../../nginx/listener/NginxEventDispatcher"
import { NginxEventListener } from "../../nginx/listener/NginxEventListener"
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import GenericNginxAction, { ActionType } from "../../nginx/actions/GenericNginxAction"
import { isAccessGranted } from "../../../core/components/accesscontrol/IsAccessGranted"
import { UserAccessLevel } from "../../user/model/UserAccessLevel"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { I18n, I18nMessage } from "../../../core/i18n/I18n"

interface NginxStatusCardState {
    loading: boolean
    running?: boolean
}

export default class NginxStatusCard extends React.Component<object, NginxStatusCardState> {
    private readonly service: NginxService
    private readonly listener: NginxEventListener

    constructor(props: object) {
        super(props)
        this.service = new NginxService()
        this.state = { loading: true }
        this.listener = () => this.handleNginxEvent()
    }

    componentDidMount() {
        NginxEventDispatcher.register(this.listener)
        this.refreshStatus()
    }

    componentWillUnmount() {
        NginxEventDispatcher.remove(this.listener)
    }

    private handleNginxEvent() {
        const { loading } = this.state
        if (loading) return

        this.setState({ loading: true }, () => this.refreshStatus())
    }

    private refreshStatus() {
        this.service
            .isRunning()
            .catch(() => undefined)
            .then(running => this.setState({ running, loading: false }))
    }

    private statusMetadata(): { color: string; label: I18nMessage } {
        const { running } = this.state

        if (running === undefined) {
            return {
                color: "var(--nginxIgnition-colorWarning)",
                label: MessageKey.FrontendHomeNginxUnknown,
            }
        }

        if (running) {
            return {
                color: "var(--nginxIgnition-colorSuccess)",
                label: MessageKey.FrontendHomeNginxOnline,
            }
        }

        return {
            color: "var(--nginxIgnition-colorError)",
            label: MessageKey.FrontendHomeNginxOffline,
        }
    }

    private confirmStop() {
        UserConfirmation.ask(MessageKey.FrontendHomeNginxStopConfirmation).then(() => {
            this.performAction(ActionType.STOP)
        })
    }

    private performAction(action: ActionType) {
        this.setState({ loading: true }, () => {
            new GenericNginxAction(action, "nginxIgnition.homeDashboard")
                .execute()
                .catch(() => {})
                .then(() => this.refreshStatus())
        })
    }

    private renderStatus() {
        const { color, label } = this.statusMetadata()

        return (
            <Flex className="home-dashboard-nginx-status" align="center">
                <span className="home-dashboard-nginx-status-dot" style={{ backgroundColor: color }} />
                <I18n id={label} />
            </Flex>
        )
    }

    private renderActions() {
        const { running } = this.state
        const readOnly = !isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.nginxServer)

        if (!running) {
            return (
                <Button
                    size="small"
                    type="primary"
                    onClick={() => this.performAction(ActionType.START)}
                    disabled={readOnly}
                >
                    <I18n id={MessageKey.FrontendHomeNginxStart} />
                </Button>
            )
        }

        return (
            <Flex className="home-dashboard-nginx-actions" gap={8} wrap="wrap">
                <Button size="small" type="primary" danger onClick={() => this.confirmStop()} disabled={readOnly}>
                    <I18n id={MessageKey.FrontendHomeNginxStop} />
                </Button>
                <Button
                    size="small"
                    type="primary"
                    onClick={() => this.performAction(ActionType.RELOAD)}
                    disabled={readOnly}
                >
                    <I18n id={MessageKey.FrontendHomeNginxReload} />
                </Button>
            </Flex>
        )
    }

    render() {
        const { loading } = this.state

        return (
            <Preloader loading={loading} size={32}>
                <div className="home-dashboard-nginx-control">
                    {this.renderStatus()}
                    <ConfigProvider componentSize="small">{this.renderActions()}</ConfigProvider>
                </div>
            </Preloader>
        )
    }
}
