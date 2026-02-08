import React from "react"
import { Flex, Statistic, Empty } from "antd"
import { Pie, Column, Area } from "@ant-design/charts"
import TrafficStatsResponse from "../model/TrafficStatsResponse"
import { formatBytes, formatNumber } from "../utils/StatsFormatters"
import {
    buildStatusDistributionData,
    buildTrafficByDomainData,
    aggregateResponses,
    STATUS_COLORS,
    buildResponseTimeData,
    buildUserAgentData,
} from "../utils/StatsChartUtils"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { I18n } from "../../../core/i18n/I18n"

interface GlobalTabProps {
    stats: TrafficStatsResponse
}

export default class GlobalTab extends React.PureComponent<GlobalTabProps> {
    private renderConnectionCards() {
        const { stats } = this.props
        return (
            <Flex className="traffic-stats-cards-row">
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficstatsConnectionsActive} />}
                        value={stats.connections.active}
                    />
                </div>
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficstatsConnectionsReading} />}
                        value={stats.connections.reading}
                    />
                </div>
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficstatsConnectionsWriting} />}
                        value={stats.connections.writing}
                    />
                </div>
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficstatsConnectionsWaiting} />}
                        value={stats.connections.waiting}
                    />
                </div>
            </Flex>
        )
    }

    private renderTotalCards() {
        const { stats } = this.props
        return (
            <Flex className="traffic-stats-cards-row">
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficstatsConnectionsAccepted} />}
                        value={formatNumber(stats.connections.accepted)}
                    />
                </div>
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficstatsConnectionsHandled} />}
                        value={formatNumber(stats.connections.handled)}
                    />
                </div>
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficstatsConnectionsRequests} />}
                        value={formatNumber(stats.connections.requests)}
                    />
                </div>
            </Flex>
        )
    }

    private renderStatusPieChart() {
        const { serverZones } = this.props.stats
        const aggregated = aggregateResponses(serverZones)
        const data = buildStatusDistributionData(aggregated)

        if (data.length === 0) {
            return <Empty description={<I18n id={MessageKey.FrontendTrafficstatsNoData} />} />
        }

        return (
            <div className="traffic-stats-chart-container">
                <p className="traffic-stats-chart-title">
                    <I18n id={MessageKey.FrontendTrafficstatsStatusDistribution} />
                </p>
                <Pie
                    data={data}
                    angleField="count"
                    colorField="status"
                    radius={0.8}
                    innerRadius={0.6}
                    label={{
                        text: "status",
                        position: "outside",
                    }}
                    legend={{
                        color: {
                            position: "bottom",
                        },
                    }}
                    scale={{
                        color: {
                            range: Object.values(STATUS_COLORS),
                        },
                    }}
                    height={300}
                />
            </div>
        )
    }

    private renderTrafficByDomainChart() {
        const { serverZones } = this.props.stats
        const data = buildTrafficByDomainData(serverZones)

        if (data.length === 0) {
            return <Empty description={<I18n id={MessageKey.FrontendTrafficstatsNoData} />} />
        }

        return (
            <div className="traffic-stats-chart-container">
                <p className="traffic-stats-chart-title">
                    <I18n id={MessageKey.FrontendTrafficstatsTrafficByDomain} />
                </p>
                <Column
                    data={data}
                    xField="name"
                    yField="requests"
                    height={300}
                    label={{
                        text: (d: { requests: number }) => formatNumber(d.requests),
                        textBaseline: "bottom",
                    }}
                    axis={{
                        x: {
                            labelAutoRotate: true,
                        },
                    }}
                />
            </div>
        )
    }

    private renderResponseTimeChart() {
        const { serverZones } = this.props.stats
        // Use global zone '*' for global response times
        const globalZone = serverZones["*"]
        if (!globalZone) return null

        const data = buildResponseTimeData(globalZone.requestMsecs)

        if (data.length === 0) {
            return null
        }

        return (
            <div className="traffic-stats-chart-container">
                <p className="traffic-stats-chart-title">
                    <I18n id={MessageKey.FrontendTrafficstatsResponseTime} />
                </p>
                <Area data={data} xField="time" yField="value" height={300} axis={{ x: { labelAutoHide: true } }} />
            </div>
        )
    }

    private renderUserAgentChart() {
        const { filterZones } = this.props.stats
        const userAgentZone = filterZones["userAgent@global"]
        if (!userAgentZone) return null

        const data = buildUserAgentData(userAgentZone)

        if (data.length === 0) {
            return null
        }

        return (
            <div className="traffic-stats-chart-container">
                <p className="traffic-stats-chart-title">
                    <I18n id={MessageKey.FrontendTrafficstatsUserAgents} />
                </p>
                <Pie
                    data={data}
                    angleField="value"
                    colorField="type"
                    radius={0.8}
                    innerRadius={0.6}
                    label={{
                        text: "type",
                        position: "outside",
                    }}
                    legend={{
                        color: {
                            position: "bottom",
                        },
                    }}
                    height={300}
                />
            </div>
        )
    }

    private renderBytesTable() {
        const { serverZones } = this.props.stats
        const data = buildTrafficByDomainData(serverZones, 10)

        return (
            <div className="traffic-stats-table-container">
                <p className="traffic-stats-chart-title">
                    <I18n id={MessageKey.FrontendTrafficstatsTrafficByDomain} />
                </p>
                <table style={{ width: "100%", borderCollapse: "collapse" }}>
                    <thead>
                        <tr style={{ borderBottom: "1px solid var(--nginxIgnition-colorBorder)" }}>
                            <th style={{ textAlign: "left", padding: "8px" }}>
                                <I18n id={MessageKey.CommonDomain} />
                            </th>
                            <th style={{ textAlign: "right", padding: "8px" }}>
                                <I18n id={MessageKey.FrontendTrafficstatsConnectionsRequests} />
                            </th>
                            <th style={{ textAlign: "right", padding: "8px" }}>
                                <I18n id={MessageKey.FrontendTrafficstatsBytesReceived} />
                            </th>
                            <th style={{ textAlign: "right", padding: "8px" }}>
                                <I18n id={MessageKey.FrontendTrafficstatsBytesSent} />
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        {data.map(item => (
                            <tr
                                key={item.name}
                                style={{ borderBottom: "1px solid var(--nginxIgnition-colorBorderSecondary)" }}
                            >
                                <td style={{ padding: "8px" }}>{item.name}</td>
                                <td style={{ textAlign: "right", padding: "8px" }}>{formatNumber(item.requests)}</td>
                                <td style={{ textAlign: "right", padding: "8px" }}>{formatBytes(item.inBytes)}</td>
                                <td style={{ textAlign: "right", padding: "8px" }}>{formatBytes(item.outBytes)}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        )
    }

    render() {
        return (
            <div className="traffic-stats-tab-content">
                {this.renderConnectionCards()}
                {this.renderTotalCards()}

                <Flex className="traffic-stats-charts-row">
                    {this.renderStatusPieChart()}
                    {this.renderTrafficByDomainChart()}
                </Flex>

                <Flex className="traffic-stats-charts-row">
                    {this.renderResponseTimeChart()}
                    {this.renderUserAgentChart()}
                </Flex>

                {this.renderBytesTable()}
            </div>
        )
    }
}
