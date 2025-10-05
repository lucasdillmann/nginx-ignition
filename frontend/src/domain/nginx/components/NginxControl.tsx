import React from "react"
import { Badge, Button, ConfigProvider, Flex, Tooltip } from "antd"
import { InfoCircleOutlined } from "@ant-design/icons"
import Preloader from "../../../core/components/preloader/Preloader"
import NginxService from "../NginxService"
import { NginxEventListener } from "../listener/NginxEventListener"
import NginxEventDispatcher from "../listener/NginxEventDispatcher"
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import "./NginxControl.css"
import GenericNginxAction, { ActionType } from "../actions/GenericNginxAction"
import { isAccessGranted } from "../../../core/components/accesscontrol/IsAccessGranted"
import { UserAccessLevel } from "../../user/model/UserAccessLevel"
import If from "../../../core/components/flowcontrol/If"
import NginxMetadata, { NginxSupportType } from "../model/NginxMetadata"

interface NginxStatusState {
    loading: boolean
    running?: boolean
    metadata?: NginxMetadata
}

export default class NginxControl extends React.Component<any, NginxStatusState> {
    private readonly service: NginxService
    private readonly listener: NginxEventListener

    constructor(props: any) {
        super(props)
        this.service = new NginxService()
        this.state = {
            loading: true,
        }
        this.listener = () => this.handleNginxEvent()
    }

    componentDidMount() {
        NginxEventDispatcher.register(this.listener)
        this.refreshNginxStatus()
        this.service
            .getMetadata()
            .then(metadata => this.setState({ metadata }))
            .catch(() => this.setState({ metadata: undefined }))
    }

    componentWillUnmount() {
        NginxEventDispatcher.remove(this.listener)
    }

    private handleNginxEvent() {
        const { loading } = this.state
        if (loading) return

        this.setState({ loading: true }, () => this.refreshNginxStatus())
    }

    private refreshNginxStatus() {
        this.service
            .isRunning()
            .catch(() => undefined)
            .then(running =>
                this.setState({
                    running: running,
                    loading: false,
                }),
            )
    }

    private renderStatusBadge(): React.ReactNode {
        const { running } = this.state

        let metadata: { color: string; description: string }
        if (running === undefined) metadata = { color: "var(--nginxIgnition-colorWarning)", description: "unknown" }
        else if (running) metadata = { color: "var(--nginxIgnition-colorSuccess)", description: "online" }
        else metadata = { color: "var(--nginxIgnition-colorError)", description: "offline" }

        return (
            <Badge
                className="nginx-control-status-badge"
                count={metadata.description}
                style={{ backgroundColor: metadata.color, borderColor: metadata.color }}
            />
        )
    }

    private confirmStop() {
        UserConfirmation.ask("Do you really want to stop the nginx server?").then(() => {
            this.performNginxAction(ActionType.STOP)
        })
    }

    private performNginxAction(action: ActionType) {
        this.setState({ loading: true }, () => {
            new GenericNginxAction(action, "nginxIgnition.nginxControl")
                .execute()
                .catch(() => {})
                .then(() => this.refreshNginxStatus())
        })
    }

    private renderActionButtons() {
        const { running } = this.state
        const readOnly = !isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.nginxServer)

        if (!running)
            return (
                <Button
                    color="primary"
                    variant="outlined"
                    onClick={() => this.performNginxAction(ActionType.START)}
                    disabled={readOnly}
                >
                    start
                </Button>
            )

        return (
            <>
                <Button color="danger" variant="outlined" onClick={() => this.confirmStop()} disabled={readOnly}>
                    stop
                </Button>
                <Button
                    className="nginx-reload-button"
                    color="primary"
                    variant="outlined"
                    onClick={() => this.performNginxAction(ActionType.RELOAD)}
                    disabled={readOnly}
                >
                    reload
                </Button>
            </>
        )
    }

    private tooltipContents() {
        const { metadata } = this.state
        if (metadata === undefined) return null

        const { version, availableSupport } = metadata

        const supportedFeatures = [
            "HTTP(S) servers",
            availableSupport.streams != NginxSupportType.NONE ? "streams" : null,
            availableSupport.tlsSni != NginxSupportType.NONE ? "TLS SNI" : null,
            availableSupport.runCode != NginxSupportType.NONE ? "JS/Lua code execution" : null,
        ].filter(feature => feature != null)

        const formatFeatures = (features: string[]): string => {
            if (features.length === 1) return features[0]
            if (features.length === 2) return `${features[0]} and ${features[1]}`

            const allButLast = features.slice(0, -1).join(", ")
            const last = features[features.length - 1]
            return `${allButLast}, and ${last}`
        }

        return `nginx version ${version} with support for ${formatFeatures(supportedFeatures)}`
    }

    render() {
        const { loading } = this.state
        const tooltipContents = this.tooltipContents()

        return (
            <Preloader loading={loading} size={32}>
                <Flex className="nginx-control-container" vertical>
                    <Flex className="nginx-control-status-title" wrap>
                        <span>server status</span>
                        <If condition={tooltipContents != undefined}>
                            <Tooltip title={tooltipContents}>
                                <InfoCircleOutlined width="10" style={{ marginLeft: 8, color: "gray" }} />
                            </Tooltip>
                        </If>
                    </Flex>
                    <Flex className="nginx-status" wrap>
                        <Flex className="nginx-status-line nginx-status-badge" align="start" justify="start">
                            {this.renderStatusBadge()}
                        </Flex>
                        <Flex className="nginx-status-line" align="end" justify="end">
                            <ConfigProvider componentSize="small">{this.renderActionButtons()}</ConfigProvider>
                        </Flex>
                    </Flex>
                </Flex>
            </Preloader>
        )
    }
}
