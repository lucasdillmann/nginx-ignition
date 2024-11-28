import React from "react";
import {Badge, Button, ConfigProvider, Flex} from "antd";
import Preloader from "../../../core/components/preloader/Preloader";
import NginxService from "../NginxService";
import Notification from "../../../core/components/notification/Notification";
import {NginxEventListener} from "../listener/NginxEventListener";
import NginxEventDispatcher from "../listener/NginxEventDispatcher";
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation";
import "./NginxControl.css"

interface NginxStatusState {
    loading: boolean,
    running?: boolean,
}

export default class NginxControl extends React.Component<any, NginxStatusState> {
    private readonly service: NginxService
    private readonly listener: NginxEventListener

    constructor(props: any) {
        super(props);
        this.service = new NginxService()
        this.state = {
            loading: true,
        }
        this.listener = () => this.handleNginxEvent()
    }

    componentDidMount() {
        NginxEventDispatcher.register(this.listener)
        this.refreshNginxStatus()
    }

    componentWillUnmount() {
        NginxEventDispatcher.remove(this.listener)
    }

    private handleNginxEvent() {
        const {loading} = this.state
        if (loading) return

        this.setState(
            { loading: true },
            () => this.refreshNginxStatus(),
        )
    }

    private refreshNginxStatus() {
        this.service
            .isRunning()
            .catch(() => undefined)
            .then(running => this.setState({
                running: running,
                loading: false,
            }))
    }

    private renderStatusBadge(): React.ReactNode {
        const {running} = this.state

        let metadata: { color: string, description: string };
        if (running === undefined)
            metadata = { color: "#a67a17", description: "unknown" }
        else if (running)
            metadata = { color: "#3f831a", description: "online" }
        else
            metadata = { color: "#af1f1f", description: "offline" }

        return (
            <Badge
                className="nginx-control-status-badge"
                count={metadata.description}
                style={{ backgroundColor: metadata.color, borderColor: metadata.color }}
            />
        )
    }

    private stopNginx() {
        UserConfirmation
            .ask("Do you really want to stop the nginx server?")
            .then(() => {
                this.performNginxAction(
                    "Stop nginx",
                    "Nginx server was stopped successfully",
                    "Nginx server failed to stop. Please check the logs for more details.",
                    () => this.service.stop(),
                )
            })
    }

    private reloadNginx() {
        this.performNginxAction(
            "Reload nginx configuration",
            "Nginx server configuration was reloaded successfully",
            "Nginx server failed to reload the configuration. Please check the logs for more details.",
            () => this.service.reloadConfiguration(),
        )
    }

    private startNginx() {
        this.performNginxAction(
            "Start nginx",
            "Nginx server was started successfully",
            "Nginx server failed to start. Please check the logs for more details.",
            () => this.service.start()
        )
    }

    private performNginxAction(
        actionName: string,
        successMessage: string,
        errorMessage: string,
        action: () => Promise<void>,
    ) {
        this.setState(
            { loading: true },
            () => {
                action()
                    .then(() => Notification.success(actionName, successMessage))
                    .catch(() => Notification.error(actionName, errorMessage))
                    .then(() => this.refreshNginxStatus())
            }
        )
    }

    private renderActionButtons() {
        const {running} = this.state

        if (!running)
            return (
                <Button
                    color="primary"
                    variant="filled"
                    onClick={() => this.startNginx()}
                >
                    start
                </Button>
            )

        return (
            <>
                <Button
                    color="danger"
                    variant="filled"
                    onClick={() => this.stopNginx()}
                >
                    stop
                </Button>
                <Button
                    className="nginx-reload-button"
                    color="primary"
                    variant="filled"
                    onClick={() => this.reloadNginx()}
                >
                    reload
                </Button>
            </>
        )
    }

    render() {
        const {loading} = this.state

        return (
            <Preloader loading={loading} size={32}>
                <Flex className="nginx-control-container" vertical>
                    <Flex className="nginx-control-status-title" wrap>
                        <span>server status</span>
                    </Flex>
                    <Flex className="nginx-status" wrap>
                        <Flex className="nginx-status-line nginx-status-badge" align="start" justify="start">
                            {this.renderStatusBadge()}
                        </Flex>
                        <Flex className="nginx-status-line" align="end" justify="end">
                            <ConfigProvider componentSize="small">
                                {this.renderActionButtons()}
                            </ConfigProvider>
                        </Flex>
                    </Flex>
                </Flex>
            </Preloader>
        )
    }
}
