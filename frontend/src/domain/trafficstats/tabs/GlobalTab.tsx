import React from "react"
import { Flex, Statistic, Empty } from "antd"
import { Column } from "@ant-design/charts"
import TrafficStatsResponse from "../model/TrafficStatsResponse"
import { formatBytes, formatNumber } from "../utils/StatsFormatters"
import {
    buildStatusDistributionData,
    buildTrafficByDomainData,
    aggregateResponses,
    buildResponseTimeData,
    buildUserAgentData,
    buildCountryCodeData,
    buildCityData,
} from "../utils/StatsChartUtils"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { I18n } from "../../../core/i18n/I18n"
import UserAgentChart from "../components/UserAgentChart"
import CountryCodeChart from "../components/CountryCodeChart"
import CityChart from "../components/CityChart"
import ResponseTimeChart from "../components/ResponseTimeChart"
import StatusDistributionChart from "../components/StatusDistributionChart"

interface GlobalTabProps {
    stats: TrafficStatsResponse
    theme: "light" | "dark"
}

export default class GlobalTab extends React.PureComponent<GlobalTabProps> {
    private renderConnectionCards() {
        const { stats } = this.props
        return (
            <Flex className="traffic-stats-cards-row">
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficStatsConnectionsActive} />}
                        value={stats.connections.active}
                    />
                </div>
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficStatsConnectionsReading} />}
                        value={stats.connections.reading}
                    />
                </div>
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficStatsConnectionsWriting} />}
                        value={stats.connections.writing}
                    />
                </div>
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficStatsConnectionsWaiting} />}
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
                        title={<I18n id={MessageKey.FrontendTrafficStatsConnectionsRequests} />}
                        value={formatNumber(stats.connections.requests)}
                    />
                </div>
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficStatsConnectionsAccepted} />}
                        value={formatNumber(stats.connections.accepted)}
                    />
                </div>
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficStatsConnectionsHandled} />}
                        value={formatNumber(stats.connections.handled)}
                    />
                </div>
            </Flex>
        )
    }

    private renderStatusPieChart() {
        const { serverZones } = this.props.stats
        const aggregated = aggregateResponses(serverZones)
        const data = buildStatusDistributionData(aggregated)
        return <StatusDistributionChart data={data} theme={this.props.theme} />
    }

    private renderTrafficByDomainChart() {
        const { serverZones } = this.props.stats
        const data = buildTrafficByDomainData(serverZones)

        if (data.length === 0) {
            return <Empty description={<I18n id={MessageKey.FrontendTrafficStatsNoData} />} />
        }

        return (
            <div className="traffic-stats-chart-container">
                <p className="traffic-stats-chart-title">
                    <I18n id={MessageKey.FrontendTrafficStatsTrafficByDomain} />
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
                    theme={this.props.theme}
                />
            </div>
        )
    }

    private renderResponseTimeChart() {
        const { serverZones } = this.props.stats
        // Use global zone '*' for global response times
        const globalZone = serverZones["*"]
        const data = globalZone ? buildResponseTimeData(globalZone.requestMsecs) : []
        return <ResponseTimeChart data={data} theme={this.props.theme} />
    }

    private renderUserAgentChart() {
        const { filterZones } = this.props.stats
        const userAgentZone = filterZones["userAgent@global"]
        const data = userAgentZone ? buildUserAgentData(userAgentZone) : []
        return <UserAgentChart data={data} theme={this.props.theme} />
    }

    private renderCountryCodeChart() {
        const { filterZones } = this.props.stats
        const countryCodeZone = filterZones["countryCode@global"]
        const data = countryCodeZone ? buildCountryCodeData(countryCodeZone) : []
        return <CountryCodeChart data={data} theme={this.props.theme} />
    }

    private renderCityChart() {
        const { filterZones } = this.props.stats
        const cityZone = filterZones["city@global"]
        const data = cityZone ? buildCityData(cityZone) : []
        return <CityChart data={data} theme={this.props.theme} />
    }

    private renderBytesTable() {
        const { serverZones } = this.props.stats
        const data = buildTrafficByDomainData(serverZones, 10)

        return (
            <div className="traffic-stats-table-container">
                <p className="traffic-stats-chart-title">
                    <I18n id={MessageKey.FrontendTrafficStatsTrafficByDomain} />
                </p>
                <table style={{ width: "100%", borderCollapse: "collapse" }}>
                    <thead>
                        <tr style={{ borderBottom: "1px solid var(--nginxIgnition-colorBorder)" }}>
                            <th style={{ textAlign: "left", padding: "8px" }}>
                                <I18n id={MessageKey.CommonDomain} />
                            </th>
                            <th style={{ textAlign: "right", padding: "8px" }}>
                                <I18n id={MessageKey.FrontendTrafficStatsConnectionsRequests} />
                            </th>
                            <th style={{ textAlign: "right", padding: "8px" }}>
                                <I18n id={MessageKey.FrontendTrafficStatsBytesReceived} />
                            </th>
                            <th style={{ textAlign: "right", padding: "8px" }}>
                                <I18n id={MessageKey.FrontendTrafficStatsBytesSent} />
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
                {this.renderTotalCards()}
                {this.renderConnectionCards()}

                <Flex className="traffic-stats-charts-row">
                    {this.renderTrafficByDomainChart()}
                    {this.renderBytesTable()}
                </Flex>

                <Flex className="traffic-stats-charts-row">{this.renderResponseTimeChart()}</Flex>

                <Flex className="traffic-stats-charts-row">
                    {this.renderStatusPieChart()}
                    {this.renderUserAgentChart()}
                </Flex>

                <Flex className="traffic-stats-charts-row">
                    {this.renderCountryCodeChart()}
                    {this.renderCityChart()}
                </Flex>
            </div>
        )
    }
}
