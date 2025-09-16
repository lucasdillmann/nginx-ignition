import React, { ReactNode } from "react"
import AccessDeniedPage from "../../core/components/accesscontrol/AccessDeniedPage"
import { UserAccessLevel } from "../user/model/UserAccessLevel"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import { Button, Flex } from "antd"
import ExportService from "./ExportService"
import AppShellContext from "../../core/components/shell/AppShellContext"
import { DatabaseOutlined, DownloadOutlined, FileZipOutlined } from "@ant-design/icons"
import Notification from "../../core/components/notification/Notification"
import { themedModal } from "../../core/components/theme/ThemedResources"
import "./ExportPage.css"

interface ExportPageState {
    nginxLoading: boolean
    databaseLoading: boolean
}

export default class ExportPage extends React.Component<any, ExportPageState> {
    private readonly service: ExportService

    constructor(props: any) {
        super(props)

        this.service = new ExportService()
        this.state = {
            nginxLoading: false,
            databaseLoading: false,
        }
    }

    private databaseBackup() {
        this.setState({ databaseLoading: true }, () =>
            this.service
                .downloadDatabaseBackup()
                .catch(error => this.showErrorNotification(error))
                .then(() => this.setState({ databaseLoading: false })),
        )
    }

    private nginxConfigurationFiles() {
        this.setState({ nginxLoading: true }, () =>
            this.service
                .downloadNginxConfigurationFiles()
                .catch(error => this.showErrorNotification(error))
                .then(() => this.setState({ nginxLoading: false })),
        )
    }

    private renderDatabaseBackup(): ReactNode {
        const { databaseLoading } = this.state
        return (
            <Flex className="export-guide-section">
                <Flex className="export-guide-section-content" vertical>
                    <Flex className="export-guide-section-title">
                        <h2>
                            <DatabaseOutlined /> Database backup
                        </h2>
                        <div className="export-guide-section-action">
                            <Button
                                type="primary"
                                size="large"
                                loading={databaseLoading}
                                onClick={() => this.databaseBackup()}
                            >
                                <DownloadOutlined /> Download
                            </Button>
                        </div>
                    </Flex>
                    <p>
                        A database backup is a snapshot of the database contents that can be used to restore the app to
                        a previous state or recover from a failure/data loss, enabling you to get nginx ignition back
                        running again. It's a good practice to backup your database regularly.
                    </p>
                    <p>
                        A backup file can take a while to be generated depending on the size of the database. The
                        produced file type will vary depending if you're using SQLite or PostgreSQL.
                    </p>
                </Flex>
            </Flex>
        )
    }

    private renderNginxConfigurationFiles(): ReactNode {
        const { nginxLoading } = this.state
        return (
            <Flex className="export-guide-section">
                <Flex className="export-guide-section-content" vertical>
                    <Flex className="export-guide-section-title">
                        <h2>
                            <FileZipOutlined /> nginx configuration files
                        </h2>
                        <div className="export-guide-section-action">
                            <Button
                                type="primary"
                                size="large"
                                loading={nginxLoading}
                                onClick={() => this.nginxConfigurationFiles()}
                            >
                                <DownloadOutlined /> Download
                            </Button>
                        </div>
                    </Flex>
                    <p>
                        Whenever you make any changes using nginx ignition and reload the server, the nginx
                        configuration files generated with all the hosts, streams, SSL certificates and more that you've
                        configured.
                    </p>
                    <p>
                        By downloading the nginx configuration files, you can analyze its contents or even deploy a
                        nginx server with these same settings and behaviours mostly the same way nginx ignition does.
                    </p>
                </Flex>
            </Flex>
        )
    }

    private showErrorNotification(error: any) {
        const onClick = () =>
            themedModal().error({
                width: 750,
                title: "Error details",
                content: <code>{error.response?.body?.message ?? error.message}</code>,
            })

        Notification.error("Download failed", "Unable to download the file at this moment. Please try again later.", {
            actions: [
                <Button key="show-details" type="default" onClick={onClick}>
                    Open error details
                </Button>,
            ],
        })
    }

    private renderPage(): ReactNode {
        return (
            <div className="export-guide-container">
                {this.renderDatabaseBackup()}
                {this.renderNginxConfigurationFiles()}
            </div>
        )
    }

    componentDidMount() {
        AppShellContext.get().updateConfig({
            title: "Export and backup",
            subtitle: "Download nginx configuration files and the ignition database contents for backup and recovery",
        })
    }

    render(): ReactNode {
        if (!isAccessGranted(UserAccessLevel.READ_ONLY, permissions => permissions.exportData)) {
            return <AccessDeniedPage />
        }

        return this.renderPage()
    }
}
