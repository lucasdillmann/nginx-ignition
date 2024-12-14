import React from "react"
import HostService from "../host/HostService"
import HostResponse from "../host/model/HostResponse"
import PaginatedSelect from "../../core/components/select/PaginatedSelect"
import { Empty, Flex, Segmented, Select } from "antd"
import {
    HddOutlined,
    ClusterOutlined,
    AuditOutlined,
    FileExcelOutlined,
    ExclamationCircleOutlined,
} from "@ant-design/icons"
import If from "../../core/components/flowcontrol/If"
import "./LogsPage.css"
import NginxService from "../nginx/NginxService"
import Preloader from "../../core/components/preloader/Preloader"
import TextArea, { TextAreaRef } from "antd/es/input/TextArea"
import AppShellContext from "../../core/components/shell/AppShellContext"
import TagGroup from "../../core/components/taggroup/TagGroup"
import SettingsDto from "../settings/model/SettingsDto"
import SettingsService from "../settings/SettingsService"
import CommonNotifications from "../../core/components/notification/CommonNotifications"
import EmptyStates from "../../core/components/emptystate/EmptyStates"

interface LogsPageState {
    settings?: SettingsDto
    hostMode: boolean
    autoRefreshSeconds?: number
    selectedHost?: HostResponse
    lineCount: number
    logType: string
    loading: boolean
    logs: string[]
    error?: Error
}

export default class LogsPage extends React.Component<any, LogsPageState> {
    private readonly hostService: HostService
    private readonly nginxService: NginxService
    private readonly settingsService: SettingsService
    private readonly contentsRef: React.RefObject<TextAreaRef>
    private refreshIntervalId?: number

    constructor(props: any) {
        super(props)
        this.hostService = new HostService()
        this.nginxService = new NginxService()
        this.settingsService = new SettingsService()
        this.contentsRef = React.createRef()
        this.state = {
            hostMode: true,
            logType: "access",
            lineCount: 50,
            loading: true,
            logs: [],
        }
    }

    componentDidMount() {
        this.settingsService
            .get()
            .then(settings => {
                this.setState({ settings })
                return this.fetchLogs()
            })
            .catch(error => {
                CommonNotifications.failedToFetch()
                this.setState({ loading: false, error })
            })

        this.configureShell()
    }

    componentDidUpdate() {
        const textarea = this.contentsRef.current?.resizableTextArea?.textArea
        if (textarea === undefined) return

        textarea.scrollTop = textarea.scrollHeight
    }

    private configureShell() {
        const { autoRefreshSeconds } = this.state
        const disabled = autoRefreshSeconds !== undefined
        const disabledReason = disabled ? "Auto refresh is enabled" : undefined

        AppShellContext.get().updateConfig({
            title: "Logs",
            subtitle: "nginx's logs for the main process or virtual hosts",
            actions: [
                {
                    description: "Refresh",
                    onClick: () => this.refreshLogs(),
                    disabled,
                    disabledReason,
                },
            ],
        })
    }

    private applyOptions() {
        const { autoRefreshSeconds } = this.state

        if (this.refreshIntervalId !== undefined) window.clearInterval(this.refreshIntervalId)

        if (autoRefreshSeconds !== undefined) {
            this.refreshIntervalId = window.setInterval(() => this.fetchLogs(), autoRefreshSeconds * 1000)
        }

        this.refreshLogs()
        this.configureShell()
    }

    private refreshLogs() {
        const { loading } = this.state
        if (loading) return

        this.setState({ loading: true }, () => this.fetchLogs())
    }

    private fetchLogs() {
        const { hostMode, lineCount, selectedHost, logType } = this.state
        if (hostMode && selectedHost === undefined) return this.setState({ loading: false })

        const logs = hostMode
            ? this.hostService.logs(selectedHost!!.id, logType, lineCount)
            : this.nginxService.logs(lineCount)

        return logs
            .then(lines => {
                this.setState({
                    loading: false,
                    logs: lines.reverse(),
                })
            })
            .catch(error => {
                CommonNotifications.failedToFetch()
                this.setState({ loading: false, error })
            })
    }

    private handleHostChange(selectedHost?: HostResponse) {
        this.setState({ selectedHost: selectedHost }, () => this.applyOptions())
    }

    private setHostMode(hostMode: boolean) {
        this.setState({ hostMode }, () => this.applyOptions())
    }

    private setLogType(logType: string) {
        this.setState({ logType }, () => this.applyOptions())
    }

    private setLineCount(lineCount: number) {
        this.setState({ lineCount }, () => this.applyOptions())
    }

    private setAutoRefreshSeconds(autoRefreshSeconds?: number) {
        this.setState({ autoRefreshSeconds }, () => this.applyOptions())
    }

    private buildLineCountOptions() {
        return [10, 25, 50, 100, 250, 500, 1000, 5000, 10000].map(item => ({
            label: item,
            value: item,
        }))
    }

    private buildAutoRefreshOptions() {
        return [1, 5, 10, 30, 60].map(item => ({
            label: `Every ${item} seconds`,
            value: item,
        }))
    }

    private handleDomainNames(domainNames?: string[]): string[] {
        if (Array.isArray(domainNames) && domainNames.length > 0) return domainNames

        return ["(default server)"]
    }

    private renderSettings() {
        const { selectedHost, hostMode, logType, lineCount, autoRefreshSeconds } = this.state
        return (
            <Flex className="log-settings-option-container">
                <Flex className="log-settings-option" vertical>
                    <p>Category</p>
                    <Segmented
                        options={[
                            { label: "Host logs", value: true, icon: <ClusterOutlined /> },
                            { label: "Server logs", value: false, icon: <HddOutlined /> },
                        ]}
                        value={hostMode}
                        onChange={value => this.setHostMode(value)}
                    />
                </Flex>
                <Flex className="log-settings-option log-settings-line-count" vertical>
                    <p>Lines</p>
                    <Select
                        options={this.buildLineCountOptions()}
                        value={lineCount}
                        onSelect={value => this.setLineCount(value)}
                    />
                </Flex>
                <Flex className="log-settings-option log-settings-auto-refresh" vertical>
                    <p>Auto-refresh</p>
                    <Select
                        placeholder="Disabled"
                        options={this.buildAutoRefreshOptions()}
                        value={autoRefreshSeconds}
                        onSelect={value => this.setAutoRefreshSeconds(value)}
                        onClear={() => this.setAutoRefreshSeconds()}
                        allowClear
                    />
                </Flex>
                <If condition={hostMode}>
                    <Flex className="log-settings-option log-settings-host" vertical>
                        <p>Host</p>
                        <PaginatedSelect
                            placeholder="Select one"
                            onChange={host => this.handleHostChange(host)}
                            pageProvider={(pageSize, pageNumber, searchTerms) =>
                                this.hostService.list(pageSize, pageNumber, searchTerms)
                            }
                            value={selectedHost}
                            itemDescription={item => (
                                <TagGroup values={this.handleDomainNames(item.domainNames)} maximumSize={1} />
                            )}
                            itemKey={item => item.id}
                            autoFocus
                        />
                    </Flex>
                    <Flex className="log-settings-option" vertical>
                        <p>Type</p>
                        <Segmented
                            options={[
                                { label: "Access logs", value: "access", icon: <AuditOutlined /> },
                                { label: "Error logs", value: "error", icon: <FileExcelOutlined /> },
                            ]}
                            value={logType}
                            onChange={value => this.setLogType(value)}
                        />
                    </Flex>
                </If>
            </Flex>
        )
    }

    private renderEmptyState(message: string, userActionOutcome: boolean = true) {
        const icon = userActionOutcome ? (
            <ExclamationCircleOutlined style={{ fontSize: 70, color: "#b8b8b8" }} />
        ) : undefined

        return <Empty image={icon} description={message} />
    }

    private renderEmptyStateIfNeeded() {
        const { settings, selectedHost, hostMode, logType, logs, loading } = this.state
        if (loading) return undefined

        const {
            nginx: { logs: logSettings },
        } = settings!!

        if (!logSettings.accessLogsEnabled && hostMode && logType === "access")
            return this.renderEmptyState("Host access logs are disabled in the nginx configuration")

        if (!logSettings.errorLogsEnabled && hostMode && logType === "error")
            return this.renderEmptyState("Host error logs are disabled in the nginx configuration")

        if (!logSettings.serverLogsEnabled && !hostMode)
            return this.renderEmptyState("nginx server logs are disabled in the nginx configuration")

        if (hostMode && selectedHost === undefined)
            return this.renderEmptyState("Please select a host in order to see its logs")

        if (logs.length === 0) return this.renderEmptyState("No logs found", false)

        return undefined
    }

    private renderLogContents() {
        const emptyState = this.renderEmptyStateIfNeeded()
        if (emptyState !== undefined) return emptyState

        const { logs } = this.state
        const contents = logs.join("\n")
        return <TextArea ref={this.contentsRef} className="log-contents-lines" value={contents} readOnly />
    }

    render() {
        const { loading, error } = this.state
        if (error !== undefined) return EmptyStates.FailedToFetch

        return (
            <Flex className="log-container" vertical>
                <Preloader loading={loading}>
                    {this.renderSettings()}

                    <Flex className="log-contents-container">{this.renderLogContents()}</Flex>
                </Preloader>
            </Flex>
        )
    }
}
