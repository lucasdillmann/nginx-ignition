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
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { I18n, i18n, I18nMessage } from "../../../core/i18n/I18n"

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

        let metadata: { color: string; description: I18nMessage }
        if (running === undefined)
            metadata = {
                color: "var(--nginxIgnition-colorWarning)",
                description: MessageKey.FrontendNginxControlStatusUnknown,
            }
        else if (running)
            metadata = {
                color: "var(--nginxIgnition-colorSuccess)",
                description: MessageKey.FrontendNginxControlStatusOnline,
            }
        else
            metadata = {
                color: "var(--nginxIgnition-colorError)",
                description: MessageKey.FrontendNginxControlStatusOffline,
            }

        return (
            <Badge
                className="nginx-control-status-badge"
                count={
                    <span>
                        <I18n id={metadata.description} />
                    </span>
                }
                style={{
                    backgroundColor: metadata.color,
                    borderColor: metadata.color,
                    padding: "2px 10px",
                    borderRadius: 20,
                }}
            />
        )
    }

    private confirmStop() {
        UserConfirmation.ask(MessageKey.FrontendNginxStopConfirmation).then(() => {
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
                    <I18n id={MessageKey.FrontendNginxControlStartButton} />
                </Button>
            )

        return (
            <>
                <Button color="danger" variant="outlined" onClick={() => this.confirmStop()} disabled={readOnly}>
                    <I18n id={MessageKey.FrontendNginxControlStopButton} />
                </Button>
                <Button
                    className="nginx-reload-button"
                    color="primary"
                    variant="outlined"
                    onClick={() => this.performNginxAction(ActionType.RELOAD)}
                    disabled={readOnly}
                >
                    <I18n id={MessageKey.FrontendNginxControlReloadButton} />
                </Button>
            </>
        )
    }

    private tooltipContents() {
        const { metadata } = this.state
        if (metadata === undefined) return null

        const { version, availableSupport } = metadata

        const supportedFeatures = [
            i18n(MessageKey.FrontendNginxControlFeatureHttp),
            availableSupport.streams != NginxSupportType.NONE
                ? i18n(MessageKey.FrontendNginxControlFeatureStreams)
                : null,
            availableSupport.tlsSni != NginxSupportType.NONE
                ? i18n(MessageKey.FrontendNginxControlFeatureTlsSni)
                : null,
            availableSupport.runCode != NginxSupportType.NONE
                ? i18n(MessageKey.FrontendNginxControlFeatureRunCode)
                : null,
        ].filter(feature => feature != null) as string[]

        const featuresList = supportedFeatures.join(", ")
        return (
            <I18n
                id={{
                    id: MessageKey.FrontendNginxControlTooltip,
                    params: { version, features: featuresList },
                }}
            />
        )
    }

    render() {
        const { loading } = this.state
        const tooltipContents = this.tooltipContents()

        return (
            <Preloader loading={loading} size={32}>
                <Flex className="nginx-control-container" vertical>
                    <Flex className="nginx-control-status-title" wrap>
                        <span>
                            <I18n id={MessageKey.FrontendNginxControlServerStatus} />
                        </span>
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
