import React from "react"
import { Tabs, Flex, Select } from "antd"
import AppShellContext from "../../core/components/shell/AppShellContext"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import { UserAccessLevel } from "../user/model/UserAccessLevel"
import AccessDeniedPage from "../../core/components/accesscontrol/AccessDeniedPage"
import Preloader from "../../core/components/preloader/Preloader"
import EmptyStates from "../../core/components/emptystate/EmptyStates"
import TrafficStatsService from "./TrafficStatsService"
import TrafficStatsResponse from "./model/TrafficStatsResponse"
import GlobalTab from "./tabs/GlobalTab"
import ByHostTab from "./tabs/ByHostTab"
import ByDomainTab from "./tabs/ByDomainTab"
import ByUpstreamTab from "./tabs/ByUpstreamTab"
import MessageKey from "../../core/i18n/model/MessageKey.generated"
import { i18n, I18n } from "../../core/i18n/I18n"
import ThemeContext from "../../core/components/context/ThemeContext"
import "./TrafficStatsPage.css"
import { Theme } from "@antv/g2/lib/spec/theme"

interface TrafficStatsPageState {
    loading: boolean
    stats?: TrafficStatsResponse
    error?: Error
    autoRefreshSeconds?: number
    activeTab: string
    theme: Theme
}

const AUTO_REFRESH_OPTIONS = [
    { label: "1s", value: 1 },
    { label: "5s", value: 5 },
    { label: "10s", value: 10 },
    { label: "30s", value: 30 },
    { label: "60s", value: 60 },
]

export default class TrafficStatsPage extends React.Component<object, TrafficStatsPageState> {
    private readonly service: TrafficStatsService
    private refreshIntervalId?: number

    constructor(props: object) {
        super(props)
        this.service = new TrafficStatsService()
        this.state = {
            loading: true,
            activeTab: "global",
            theme: ThemeContext.isDarkMode() ? "dark" : "light",
        }
    }

    componentDidMount() {
        this.configureShell()
        this.fetchStats()
        ThemeContext.register(this.handleThemeChange.bind(this))
    }

    componentWillUnmount() {
        this.stopAutoRefresh()
        ThemeContext.deregister(this.handleThemeChange.bind(this))
    }

    private handleThemeChange(darkMode: boolean) {
        this.setState({ theme: darkMode ? "dark" : "light" })
    }

    private configureShell() {
        const { autoRefreshSeconds } = this.state
        const autoRefreshEnabled = autoRefreshSeconds !== undefined && autoRefreshSeconds > 0

        AppShellContext.get().updateConfig({
            title: MessageKey.CommonTrafficStats,
            subtitle: MessageKey.FrontendTrafficStatsSubtitle,
            actions: [
                {
                    description: MessageKey.CommonRefresh,
                    onClick: () => this.fetchStats(),
                    disabled: autoRefreshEnabled,
                    disabledReason: autoRefreshEnabled
                        ? MessageKey.FrontendTrafficStatsAutoRefreshDisabledReason
                        : undefined,
                },
            ],
        })
    }

    private async fetchStats() {
        this.setState({ loading: true })
        try {
            const stats = await this.service.getStats()
            this.setState({ stats, loading: false, error: undefined })
        } catch (error) {
            this.setState({ loading: false, error: error as Error })
        }
    }

    private stopAutoRefresh() {
        if (this.refreshIntervalId !== undefined) {
            window.clearInterval(this.refreshIntervalId)
            this.refreshIntervalId = undefined
            this.setState({ autoRefreshSeconds: undefined })
        }
    }

    private applyAutoRefresh(seconds: number) {
        this.stopAutoRefresh()

        const autoRefreshSeconds = seconds > 0 ? seconds : undefined
        this.setState({ autoRefreshSeconds }, () => {
            this.configureShell()
            if (autoRefreshSeconds !== undefined) {
                this.refreshIntervalId = window.setInterval(() => this.fetchStats(), autoRefreshSeconds * 1000)
            }
        })
    }

    private renderSettings() {
        const { autoRefreshSeconds } = this.state

        return (
            <div className="traffic-stats-settings-wrapper">
                <div className="traffic-stats-settings-option">
                    <p>
                        <I18n id={MessageKey.CommonAutoRefresh} />
                    </p>
                    <Select
                        className="traffic-stats-auto-refresh"
                        placeholder={i18n(MessageKey.CommonDisabled)}
                        options={AUTO_REFRESH_OPTIONS}
                        value={autoRefreshSeconds}
                        onSelect={value => this.applyAutoRefresh(value)}
                        onClear={() => this.stopAutoRefresh()}
                        allowClear
                    />
                </div>
            </div>
        )
    }

    private renderTabs() {
        const { stats, activeTab, theme } = this.state

        if (!stats) return null

        const items = [
            {
                key: "global",
                label: <I18n id={MessageKey.FrontendTrafficStatsGlobalTab} />,
                children: <GlobalTab stats={stats} theme={theme} />,
            },
            {
                key: "byHost",
                label: <I18n id={MessageKey.FrontendTrafficStatsByHostTab} />,
                children: <ByHostTab stats={stats} theme={theme} />,
            },
            {
                key: "byDomain",
                label: <I18n id={MessageKey.FrontendTrafficStatsByDomainTab} />,
                children: <ByDomainTab stats={stats} theme={theme} />,
            },
            {
                key: "byUpstream",
                label: <I18n id={MessageKey.FrontendTrafficStatsByUpstreamTab} />,
                children: <ByUpstreamTab stats={stats} theme={theme} />,
            },
        ]

        return (
            <div className="traffic-stats-tabs-container">
                <Tabs
                    activeKey={activeTab}
                    items={items}
                    onChange={key => this.setState({ activeTab: key })}
                    tabBarExtraContent={{ right: this.renderSettings() }}
                />
            </div>
        )
    }

    render() {
        if (!isAccessGranted(UserAccessLevel.READ_ONLY, permissions => permissions.trafficStats)) {
            return <AccessDeniedPage />
        }

        const { loading, error, stats } = this.state

        return (
            <Flex className="traffic-stats-container" vertical>
                <Preloader loading={loading}>
                    {error && EmptyStates.FailedToFetch}
                    {!error && stats && this.renderTabs()}
                </Preloader>
            </Flex>
        )
    }
}
