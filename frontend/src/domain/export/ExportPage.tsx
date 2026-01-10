import React, { ReactNode } from "react"
import AccessDeniedPage from "../../core/components/accesscontrol/AccessDeniedPage"
import { UserAccessLevel } from "../user/model/UserAccessLevel"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import { Button, Flex, Form, Input, Modal, Space } from "antd"
import ExportService from "./ExportService"
import AppShellContext from "../../core/components/shell/AppShellContext"
import { DatabaseOutlined, DownloadOutlined, FileZipOutlined, QuestionCircleOutlined } from "@ant-design/icons"
import Notification from "../../core/components/notification/Notification"
import { themedModal } from "../../core/components/theme/ThemedResources"
import "./ExportPage.css"

interface ExportPageState {
    nginxModalOpen: boolean
    nginxBasePath: string
    nginxConfigPath: string
    nginxLogPath: string
    nginxCachePath: string
    nginxTempPath: string
    nginxLoading: boolean
    databaseLoading: boolean
}

export default class ExportPage extends React.Component<any, ExportPageState> {
    private readonly service: ExportService

    constructor(props: any) {
        super(props)

        this.service = new ExportService()
        this.state = {
            nginxModalOpen: false,
            nginxLoading: false,
            nginxBasePath: "",
            nginxConfigPath: "",
            nginxLogPath: "",
            nginxCachePath: "",
            nginxTempPath: "",
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

    private openDatabaseHelpGuide() {
        window.open(
            "https://github.com/lucasdillmann/nginx-ignition/blob/main/docs/database-restore-guide.md",
            "_blank",
            "noopener",
        )
    }

    private openNginxModal() {
        this.setState({ nginxModalOpen: true })
    }

    private closeNginxModal() {
        this.setState({ nginxModalOpen: false })
    }

    private nginxConfigurationFiles() {
        const { nginxBasePath, nginxConfigPath, nginxLogPath, nginxCachePath, nginxTempPath } = this.state

        this.setState({ nginxLoading: true, nginxModalOpen: false }, () =>
            this.service
                .downloadNginxConfigurationFiles(
                    nginxBasePath,
                    nginxConfigPath,
                    nginxLogPath,
                    nginxCachePath,
                    nginxTempPath,
                )
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
                                type="default"
                                size="large"
                                onClick={() => this.openDatabaseHelpGuide()}
                                style={{ marginRight: 10 }}
                            >
                                <QuestionCircleOutlined />
                            </Button>
                            <Space />
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
                        A database backup is a snapshot of the nginx ignition configurations (hosts, certificates,
                        streams, etc) that can be used to restore the app to a previous state or recover from a
                        failure/data loss, enabling you to get nginx ignition back and running again.
                    </p>
                    <p>
                        It's a good idea to backup your database regularly using an automated process whenever possible
                        (which can be done using cron jobs or similar tools calling the nginx ignition API to generate
                        the backup file or by interacting with the database directly).
                    </p>
                    <p>
                        Please note that a backup file can take a while to be generated depending on the size of the
                        database. The produced file type will vary depending if you're using SQLite or PostgreSQL (more
                        details available at the nginx ignition's documentations).
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
                            <FileZipOutlined /> nginx configuration
                        </h2>
                        <div className="export-guide-section-action">
                            <Button
                                type="primary"
                                size="large"
                                loading={nginxLoading}
                                onClick={() => this.openNginxModal()}
                            >
                                <DownloadOutlined /> Download
                            </Button>
                        </div>
                    </Flex>
                    <p>
                        Whenever you make a change using nginx ignition and reload the server, the nginx configuration
                        files are generated or updated with all the hosts, streams, SSL certificates and more that
                        you've enabled.
                    </p>
                    <p>
                        By downloading the nginx configuration files, you can analyze/review its contents or even deploy
                        a nginx server with these same settings and behaviours mostly the same way nginx ignition does.
                    </p>
                </Flex>
            </Flex>
        )
    }

    private setValue(field: string, value: string) {
        this.setState(current => ({
            ...current,
            [field]: value,
        }))
    }

    private renderNginxConfigurationModal(): ReactNode {
        const {
            nginxModalOpen,
            nginxBasePath,
            nginxConfigPath,
            nginxLogPath,
            nginxCachePath,
            nginxTempPath,
            nginxLoading,
        } = this.state
        return (
            <Modal
                afterClose={() => this.closeNginxModal()}
                onCancel={() => this.closeNginxModal()}
                onOk={() => this.nginxConfigurationFiles()}
                okButtonProps={{
                    disabled: nginxLoading,
                }}
                title="nginx configuration"
                width={800}
                open={nginxModalOpen}
                okText="Continue"
                cancelText="Cancel"
            >
                <p>
                    If needed, you can customize the paths that the configuration files should use by filling the fields
                    below. This action is optional and, if left empty, the files will be generated using relative paths.
                </p>
                <br />
                <Form.Item label="Base path for the nginx files" initialValue={nginxBasePath} layout="vertical">
                    <Input
                        size="large"
                        onChange={event => this.setValue("nginxBasePath", event.target.value)}
                        required={false}
                        placeholder="e.g. /etc/nginx"
                        autoFocus
                    />
                </Form.Item>
                <Form.Item
                    label="Path for the nginx configuration files"
                    initialValue={nginxConfigPath}
                    layout="vertical"
                >
                    <Input
                        size="large"
                        onChange={event => this.setValue("nginxConfigPath", event.target.value)}
                        required={false}
                        placeholder="e.g. /etc/nginx/config"
                        autoFocus
                    />
                </Form.Item>
                <Form.Item label="Path for the nginx log files" initialValue={nginxLogPath} layout="vertical">
                    <Input
                        size="large"
                        onChange={event => this.setValue("nginxLogPath", event.target.value)}
                        required={false}
                        placeholder="e.g. /var/log/nginx"
                        autoFocus
                    />
                </Form.Item>
                <Form.Item label="Path for the nginx cache files" initialValue={nginxCachePath} layout="vertical">
                    <Input
                        size="large"
                        onChange={event => this.setValue("nginxCachePath", event.target.value)}
                        required={false}
                        placeholder="e.g. /var/run/nginx/cache"
                        autoFocus
                    />
                </Form.Item>
                <Form.Item label="Path for the nginx temp files" initialValue={nginxTempPath} layout="vertical">
                    <Input
                        size="large"
                        onChange={event => this.setValue("nginxTempPath", event.target.value)}
                        required={false}
                        placeholder="e.g. /tmp/nginx"
                        autoFocus
                    />
                </Form.Item>
            </Modal>
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
                {this.renderNginxConfigurationModal()}
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
