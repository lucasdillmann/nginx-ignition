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
import MessageKey from "../../core/i18n/model/MessageKey.generated"
import { I18n, i18n } from "../../core/i18n/I18n"

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
                            <DatabaseOutlined /> <I18n id={MessageKey.FrontendExportSectionDatabaseTitle} />
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
                                <DownloadOutlined /> <I18n id={MessageKey.CommonDownload} />
                            </Button>
                        </div>
                    </Flex>
                    <p>
                        <I18n id={MessageKey.FrontendExportSectionDatabaseDescription1} />
                    </p>
                    <p>
                        <I18n id={MessageKey.FrontendExportSectionDatabaseDescription2} />
                    </p>
                    <p>
                        <I18n id={MessageKey.FrontendExportSectionDatabaseDescription3} />
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
                            <FileZipOutlined /> <I18n id={MessageKey.CommonNginxConfiguration} />
                        </h2>
                        <div className="export-guide-section-action">
                            <Button
                                type="primary"
                                size="large"
                                loading={nginxLoading}
                                onClick={() => this.openNginxModal()}
                            >
                                <DownloadOutlined /> <I18n id={MessageKey.CommonDownload} />
                            </Button>
                        </div>
                    </Flex>
                    <p>
                        <I18n id={MessageKey.FrontendExportSectionNginxDescription1} />
                    </p>
                    <p>
                        <I18n id={MessageKey.FrontendExportSectionNginxDescription2} />
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
                title={<I18n id={MessageKey.CommonNginxConfiguration} />}
                width={800}
                open={nginxModalOpen}
                okText={<I18n id={MessageKey.CommonContinue} />}
                cancelText={<I18n id={MessageKey.CommonCancel} />}
            >
                <p>
                    <I18n id={MessageKey.FrontendExportNginxModalDescription} />
                </p>
                <br />
                <Form.Item
                    label={<I18n id={MessageKey.FrontendExportNginxModalBasePath} />}
                    initialValue={nginxBasePath}
                    layout="vertical"
                >
                    <Input
                        size="large"
                        onChange={event => this.setValue("nginxBasePath", event.target.value)}
                        required={false}
                        placeholder={i18n(MessageKey.FrontendExportNginxModalBasePathPlaceholder)}
                        autoFocus
                    />
                </Form.Item>
                <Form.Item
                    label={<I18n id={MessageKey.FrontendExportNginxModalConfigPath} />}
                    initialValue={nginxConfigPath}
                    layout="vertical"
                >
                    <Input
                        size="large"
                        onChange={event => this.setValue("nginxConfigPath", event.target.value)}
                        required={false}
                        placeholder={i18n(MessageKey.FrontendExportNginxModalConfigPathPlaceholder)}
                        autoFocus
                    />
                </Form.Item>
                <Form.Item
                    label={<I18n id={MessageKey.FrontendExportNginxModalLogPath} />}
                    initialValue={nginxLogPath}
                    layout="vertical"
                >
                    <Input
                        size="large"
                        onChange={event => this.setValue("nginxLogPath", event.target.value)}
                        required={false}
                        placeholder={i18n(MessageKey.FrontendExportNginxModalLogPathPlaceholder)}
                        autoFocus
                    />
                </Form.Item>
                <Form.Item
                    label={<I18n id={MessageKey.FrontendExportNginxModalCachePath} />}
                    initialValue={nginxCachePath}
                    layout="vertical"
                >
                    <Input
                        size="large"
                        onChange={event => this.setValue("nginxCachePath", event.target.value)}
                        required={false}
                        placeholder={i18n(MessageKey.FrontendExportNginxModalCachePathPlaceholder)}
                        autoFocus
                    />
                </Form.Item>
                <Form.Item
                    label={<I18n id={MessageKey.FrontendExportNginxModalTempPath} />}
                    initialValue={nginxTempPath}
                    layout="vertical"
                >
                    <Input
                        size="large"
                        onChange={event => this.setValue("nginxTempPath", event.target.value)}
                        required={false}
                        placeholder={i18n(MessageKey.FrontendExportNginxModalTempPathPlaceholder)}
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
                title: <I18n id={MessageKey.FrontendComponentsErrorDetails} />,
                content: <code>{error.response?.body?.message ?? error.message}</code>,
            })

        Notification.error(
            MessageKey.FrontendExportDownloadFailed,
            MessageKey.FrontendExportDownloadFailedDescription,
            {
                actions: [
                    <Button key="show-details" type="default" onClick={onClick}>
                        Open error details
                    </Button>,
                ],
            },
        )
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
            title: MessageKey.CommonExportAndBackup,
            subtitle: MessageKey.FrontendExportSubtitle,
        })
    }

    render(): ReactNode {
        if (!isAccessGranted(UserAccessLevel.READ_ONLY, permissions => permissions.exportData)) {
            return <AccessDeniedPage />
        }

        return this.renderPage()
    }
}
