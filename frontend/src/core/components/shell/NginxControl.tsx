import React from "react";
import {Badge, Button, ConfigProvider, Flex} from "antd";
import Preloader from "../preloader/Preloader";
import NginxService from "../../../domain/nginx/NginxService";
import NotificationFacade from "../notification/NotificationFacade";
import styles from "./NginxControl.styles"
import {NginxEventListener} from "../../../domain/nginx/listener/NginxEventListener";
import NginxEventDispatcher from "../../../domain/nginx/listener/NginxEventDispatcher";

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
                className="site-badge-count-109"
                count={metadata.description}
                style={{ backgroundColor: metadata.color, borderColor: metadata.color }}
            />
        )
    }

    private stopNginx() {
        this.performNginxAction(
            "Stop nginx",
            "Nginx server was stopped successfully",
            "Nginx server failed to stop. Please check the logs for more details.",
            () => this.service.stop(),
        )
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
                    .then(() => NotificationFacade.success(actionName, successMessage))
                    .catch(() => NotificationFacade.error(actionName, errorMessage))
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
                    style={{ marginLeft: 5 }}
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

        const components = (
            <Flex style={styles.container} vertical>
                <Flex wrap>
                    <span>nginx status</span>
                </Flex>
                <Flex style={{ marginTop: 15 }} wrap>
                    <Flex style={{ width: "50%", height: 20 }} align="start" justify="start">
                        {this.renderStatusBadge()}
                    </Flex>
                    <Flex style={{ width: "50%", height: 20 }} align="end" justify="end">
                        <ConfigProvider componentSize="small">
                            {this.renderActionButtons()}
                        </ConfigProvider>
                    </Flex>
                </Flex>
            </Flex>
        )

        return loading ? <Preloader size={32}>{components}</Preloader> : components
    }
}