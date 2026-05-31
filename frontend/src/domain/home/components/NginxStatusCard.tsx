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
import { I18n, I18nMessage, i18n } from "../../../core/i18n/I18n"
import If from "../../../core/components/flowcontrol/If"

interface NginxStatusCardProps {
    refreshToken: number
}

interface NginxStatusCardState {
    loading: boolean
    running?: boolean
    uptimeSeconds?: number
}

export default class NginxStatusCard extends React.Component<NginxStatusCardProps, NginxStatusCardState> {
    private readonly service: NginxService
    private readonly listener: NginxEventListener

    constructor(props: NginxStatusCardProps) {
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

    componentDidUpdate(previousProps: NginxStatusCardProps) {
        if (previousProps.refreshToken !== this.props.refreshToken) {
            this.refreshStatus()
        }
    }

    private handleNginxEvent() {
        const { loading } = this.state
        if (loading) return

        this.setState({ loading: true }, () => this.refreshStatus())
    }

    private refreshStatus() {
        this.service
            .getStatus()
            .catch(() => undefined)
            .then(status =>
                this.setState({
                    running: status?.running,
                    uptimeSeconds: status?.uptimeSeconds,
                    loading: false,
                }),
            )
    }

    private formatCountUnit(count: number, singularKey: I18nMessage, pluralKey: I18nMessage): string {
        const unit = i18n(count === 1 ? singularKey : pluralKey)
        return `${count} ${unit}`
    }

    private uptimeLabelParams(totalSeconds: number): { days: string; hours: string; minutes: string; seconds: string } {
        const dayCount = Math.floor(totalSeconds / 86400)
        const hourCount = Math.floor((totalSeconds % 86400) / 3600)
        const minuteCount = Math.floor((totalSeconds % 3600) / 60)
        const secondCount = totalSeconds % 60

        return {
            days: this.formatCountUnit(dayCount, MessageKey.CommonTimeUnitDay, MessageKey.CommonTimeUnitDays),
            hours: this.formatCountUnit(hourCount, MessageKey.CommonTimeUnitHour, MessageKey.CommonTimeUnitHours),
            minutes: this.formatCountUnit(
                minuteCount,
                MessageKey.CommonTimeUnitMinute,
                MessageKey.CommonTimeUnitMinutes,
            ),
            seconds: this.formatCountUnit(secondCount, MessageKey.CommonTimeUnitSecond, MessageKey.CommonUnitSeconds),
        }
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
        const { running, uptimeSeconds } = this.state
        const { color, label } = this.statusMetadata()
        const showUptime = running === true && Boolean(uptimeSeconds)

        return (
            <Flex className="home-dashboard-nginx-status" align="stretch">
                <span className="home-dashboard-nginx-status-border" style={{ backgroundColor: color }} />
                <Flex className="home-dashboard-nginx-status-label" align="center">
                    <If condition={showUptime}>
                        <I18n
                            id={{
                                id: MessageKey.FrontendHomeNginxOnlineUptime,
                                params: this.uptimeLabelParams(uptimeSeconds!),
                            }}
                        />
                    </If>
                    <If condition={!showUptime}>
                        <I18n id={label} />
                    </If>
                </Flex>
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
