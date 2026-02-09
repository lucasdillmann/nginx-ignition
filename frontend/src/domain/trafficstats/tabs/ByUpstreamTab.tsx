import React from "react"
import { Flex, Select, Statistic, Empty, Table, Tag } from "antd"
import { Pie, Area } from "@ant-design/charts"
import TrafficStatsResponse, { UpstreamZoneData } from "../model/TrafficStatsResponse"
import { formatBytes, formatNumber, formatMs } from "../utils/StatsFormatters"
import { STATUS_COLORS, StatusDataItem } from "../utils/StatsChartUtils"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { I18n } from "../../../core/i18n/I18n"
import { CheckCircleOutlined, CloseCircleOutlined, WarningOutlined } from "@ant-design/icons"
import StatusDistributionChart from "../components/StatusDistributionChart"

interface ByUpstreamTabProps {
    stats: TrafficStatsResponse
    theme: "light" | "dark"
    selectedUpstream?: string
    onSelectUpstream: (upstream: string) => void
}

export default class ByUpstreamTab extends React.Component<ByUpstreamTabProps> {
    private getSelectedUpstreamData(): UpstreamZoneData[] | undefined {
        const { upstreamZones } = this.props.stats
        const { selectedUpstream } = this.props
        if (!selectedUpstream) return undefined
        return upstreamZones[selectedUpstream]
    }

    private getTotalStats(servers: UpstreamZoneData[]) {
        return servers.reduce(
            (acc, server) => ({
                requests: acc.requests + server.requestCounter,
                inBytes: acc.inBytes + server.inBytes,
                outBytes: acc.outBytes + server.outBytes,
            }),
            { requests: 0, inBytes: 0, outBytes: 0 },
        )
    }

    private getAvgResponseTime(servers: UpstreamZoneData[]): number {
        const totalRequests = servers.reduce((sum, s) => sum + s.requestCounter, 0)
        if (totalRequests === 0) return 0
        const totalMs = servers.reduce((sum, s) => sum + s.responseMsecCounter, 0)
        return totalMs / totalRequests
    }

    private renderUpstreamSelector() {
        const { upstreamZones } = this.props.stats
        const { selectedUpstream, onSelectUpstream } = this.props
        const upstreams = Object.keys(upstreamZones)

        const options = upstreams.map(name => ({
            value: name,
            label: name,
        }))

        return (
            <Flex className="traffic-stats-settings-option">
                <p>
                    <I18n id={MessageKey.FrontendTrafficStatsUpstreamServers} />
                </p>
                <Select
                    className="traffic-stats-selector"
                    placeholder={<I18n id={MessageKey.FrontendTrafficStatsSelectUpstream} />}
                    options={options}
                    value={selectedUpstream}
                    onChange={value => onSelectUpstream(value)}
                    showSearch
                    filterOption={(input, option) =>
                        (option?.label?.toString() ?? "").toLowerCase().includes(input.toLowerCase())
                    }
                />
            </Flex>
        )
    }

    private renderStatCards(servers: UpstreamZoneData[]) {
        const totals = this.getTotalStats(servers)
        const avgResponseTime = this.getAvgResponseTime(servers)
        const upServers = servers.filter(s => !s.down).length

        return (
            <Flex className="traffic-stats-cards-row">
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficStatsUpstreamServers} />}
                        value={`${upServers}/${servers.length}`}
                    />
                </div>
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficStatsConnectionsRequests} />}
                        value={formatNumber(totals.requests)}
                    />
                </div>
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficStatsBytesReceived} />}
                        value={formatBytes(totals.inBytes)}
                    />
                </div>
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficStatsAverageResponseTime} />}
                        value={formatMs(avgResponseTime)}
                    />
                </div>
            </Flex>
        )
    }

    private renderServerTable(servers: UpstreamZoneData[]) {
        const columns = [
            {
                title: <I18n id={MessageKey.FrontendTrafficStatsServer} />,
                dataIndex: "server",
                key: "server",
            },
            {
                title: <I18n id={MessageKey.FrontendTrafficStatsServerStatus} />,
                key: "status",
                render: (_: unknown, record: UpstreamZoneData) => {
                    if (record.down) {
                        return (
                            <Tag color="error" icon={<CloseCircleOutlined />}>
                                <I18n id={MessageKey.FrontendTrafficStatsUpstreamDown} />
                            </Tag>
                        )
                    }
                    if (record.backup) {
                        return (
                            <Tag color="warning" icon={<WarningOutlined />}>
                                <I18n id={MessageKey.FrontendTrafficStatsUpstreamBackup} />
                            </Tag>
                        )
                    }
                    return (
                        <Tag color="success" icon={<CheckCircleOutlined />}>
                            <I18n id={MessageKey.FrontendTrafficStatsConnectionsActive} />
                        </Tag>
                    )
                },
            },
            {
                title: <I18n id={MessageKey.FrontendStreamComponentsBackendsettingsWeight} />,
                dataIndex: "weight",
                key: "weight",
            },
            {
                title: <I18n id={MessageKey.FrontendTrafficStatsConnectionsRequests} />,
                key: "requests",
                render: (_: unknown, record: UpstreamZoneData) => formatNumber(record.requestCounter),
            },
            {
                title: <I18n id={MessageKey.FrontendTrafficStatsAverageResponseTime} />,
                key: "responseTime",
                render: (_: unknown, record: UpstreamZoneData) => {
                    if (record.requestCounter === 0) return "-"
                    return formatMs(record.responseMsecCounter / record.requestCounter)
                },
            },
            {
                title: <I18n id={MessageKey.FrontendTrafficStatsBytesSent} />,
                key: "outBytes",
                render: (_: unknown, record: UpstreamZoneData) => formatBytes(record.outBytes),
            },
        ]

        return (
            <div className="traffic-stats-table-container">
                <p className="traffic-stats-chart-title">
                    <I18n id={MessageKey.FrontendTrafficStatsUpstreamServers} />
                </p>
                <Table dataSource={servers} columns={columns} pagination={false} rowKey="server" size="small" />
            </div>
        )
    }

    private renderResponseTimeChart(servers: UpstreamZoneData[]) {
        const data: { time: string; timestamp: number; value: number; server: string }[] = []

        servers.forEach(server => {
            if (!server?.requestMsecs?.times) return

            server.requestMsecs.times.forEach((t, i) => {
                if (server.requestMsecs.msecs[i] > 0) {
                    data.push({
                        time: new Date(t).toLocaleTimeString(),
                        timestamp: t,
                        value: server.requestMsecs.msecs[i],
                        server: server.server,
                    })
                }
            })
        })

        if (data.length === 0) return null

        data.sort((a, b) => a.timestamp - b.timestamp)

        return (
            <div className="traffic-stats-chart-container">
                <p className="traffic-stats-chart-title">
                    <I18n id={MessageKey.FrontendTrafficStatsResponseTime} />
                </p>
                <Area
                    data={data}
                    xField="time"
                    yField="value"
                    seriesField="server"
                    height={300}
                    axis={{ x: { labelAutoHide: true } }}
                    legend={{ position: "bottom" }}
                    theme={this.props.theme}
                />
            </div>
        )
    }

    private renderDistributionPieChart(servers: UpstreamZoneData[]) {
        const data = servers
            .filter(s => s.requestCounter > 0)
            .map(s => ({
                server: s.server,
                requests: s.requestCounter,
            }))

        if (data.length === 0) {
            return null
        }

        return (
            <div className="traffic-stats-chart-container">
                <p className="traffic-stats-chart-title">
                    <I18n id={MessageKey.FrontendTrafficStatsTrafficDistribution} />
                </p>
                <Pie
                    data={data}
                    angleField="requests"
                    colorField="server"
                    radius={0.8}
                    innerRadius={0.6}
                    label={{
                        text: "server",
                        position: "outside",
                    }}
                    legend={{
                        color: {
                            position: "bottom",
                        },
                    }}
                    height={300}
                    theme={this.props.theme}
                />
            </div>
        )
    }

    private renderStatusDistributionChart(servers: UpstreamZoneData[]) {
        const aggregated = servers.reduce(
            (acc, server) => ({
                "1xx": acc["1xx"] + server.responses["1xx"],
                "2xx": acc["2xx"] + server.responses["2xx"],
                "3xx": acc["3xx"] + server.responses["3xx"],
                "4xx": acc["4xx"] + server.responses["4xx"],
                "5xx": acc["5xx"] + server.responses["5xx"],
            }),
            { "1xx": 0, "2xx": 0, "3xx": 0, "4xx": 0, "5xx": 0 },
        )

        const data: StatusDataItem[] = [
            { status: "1xx", count: aggregated["1xx"], color: STATUS_COLORS["1xx"] },
            { status: "2xx", count: aggregated["2xx"], color: STATUS_COLORS["2xx"] },
            { status: "3xx", count: aggregated["3xx"], color: STATUS_COLORS["3xx"] },
            { status: "4xx", count: aggregated["4xx"], color: STATUS_COLORS["4xx"] },
            { status: "5xx", count: aggregated["5xx"], color: STATUS_COLORS["5xx"] },
        ].filter(item => item.count > 0)

        if (data.length === 0) {
            return null
        }

        return <StatusDistributionChart data={data} theme={this.props.theme} />
    }

    render() {
        const servers = this.getSelectedUpstreamData()

        return (
            <div className="traffic-stats-tab-content">
                {this.renderUpstreamSelector()}
                {servers ? (
                    <>
                        {this.renderStatCards(servers)}
                        <Flex className="traffic-stats-charts-row">{this.renderResponseTimeChart(servers)}</Flex>
                        {this.renderServerTable(servers)}
                        <Flex className="traffic-stats-charts-row">
                            {this.renderDistributionPieChart(servers)}
                            {this.renderStatusDistributionChart(servers)}
                        </Flex>
                    </>
                ) : (
                    <Empty description={<I18n id={MessageKey.FrontendTrafficStatsSelectUpstream} />} />
                )}
            </div>
        )
    }
}
