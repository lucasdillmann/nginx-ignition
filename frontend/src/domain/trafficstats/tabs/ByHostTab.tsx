import React from "react"
import { Flex, Select, Statistic, Empty, Table } from "antd"
import { Pie, Area } from "@ant-design/charts"
import TrafficStatsResponse, { ZoneData } from "../model/TrafficStatsResponse"
import HostService from "../../host/HostService"
import HostResponse from "../../host/model/HostResponse"
import { formatBytes, formatNumber, formatMs } from "../utils/StatsFormatters"
import {
    buildStatusDistributionData,
    STATUS_COLORS,
    buildResponseTimeData,
    buildUserAgentData,
    buildCountryCodeData,
} from "../utils/StatsChartUtils"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { I18n, i18n } from "../../../core/i18n/I18n"

interface ByHostTabProps {
    stats: TrafficStatsResponse
    theme: "light" | "dark"
}

interface ByHostTabState {
    hosts: HostResponse[]
    selectedHostId?: string
    loading: boolean
}

export default class ByHostTab extends React.Component<ByHostTabProps, ByHostTabState> {
    private readonly hostService: HostService

    constructor(props: ByHostTabProps) {
        super(props)
        this.hostService = new HostService()
        this.state = {
            hosts: [],
            loading: true,
        }
    }

    componentDidMount() {
        this.loadHosts()
    }

    private async loadHosts() {
        try {
            const page = await this.hostService.list(100, 0)
            this.setState({ hosts: page.contents, loading: false })
        } catch {
            this.setState({ loading: false })
        }
    }

    private getHostLabel(hostId: string): string {
        const host = this.state.hosts.find(h => h.id === hostId)
        if (host) {
            if (host.defaultServer) return i18n(MessageKey.CommonDefaultServerLabel)
            if (host.domainNames && host.domainNames.length > 0) return host.domainNames.join(", ")
        }
        return hostId
    }

    private getSelectedZoneData(): ZoneData | undefined {
        const { filterZones } = this.props.stats
        const { selectedHostId } = this.state
        if (!selectedHostId || !filterZones.hosts) return undefined
        return filterZones.hosts[selectedHostId]
    }

    private getAvgResponseTime(zone: ZoneData): number {
        if (zone.requestCounter === 0) return 0
        return zone.requestMsecCounter / zone.requestCounter
    }

    private renderHostSelector() {
        const { filterZones } = this.props.stats
        const { selectedHostId, loading } = this.state
        const hostIds = filterZones.hosts ? Object.keys(filterZones.hosts) : []

        const options = hostIds.map(id => ({
            value: id,
            label: this.getHostLabel(id),
        }))

        return (
            <Flex className="traffic-stats-settings-option">
                <p>
                    <I18n id={MessageKey.CommonHost} />
                </p>
                <Select
                    className="traffic-stats-selector"
                    placeholder={<I18n id={MessageKey.FrontendTrafficStatsSelectHost} />}
                    options={options}
                    value={selectedHostId}
                    onChange={value => this.setState({ selectedHostId: value })}
                    loading={loading}
                    showSearch
                    filterOption={(input, option) =>
                        (option?.label?.toString() ?? "").toLowerCase().includes(input.toLowerCase())
                    }
                />
            </Flex>
        )
    }

    private renderStatCards(zone: ZoneData) {
        return (
            <Flex className="traffic-stats-cards-row">
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficStatsConnectionsRequests} />}
                        value={formatNumber(zone.requestCounter)}
                    />
                </div>
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficStatsBytesReceived} />}
                        value={formatBytes(zone.inBytes)}
                    />
                </div>
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficStatsBytesSent} />}
                        value={formatBytes(zone.outBytes)}
                    />
                </div>
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficStatsAverageResponseTime} />}
                        value={formatMs(this.getAvgResponseTime(zone))}
                    />
                </div>
            </Flex>
        )
    }

    private renderStatusPieChart(zone: ZoneData) {
        const data = buildStatusDistributionData(zone.responses)

        if (data.length === 0) {
            return <Empty description={<I18n id={MessageKey.FrontendTrafficStatsNoData} />} />
        }

        return (
            <div className="traffic-stats-chart-container">
                <p className="traffic-stats-chart-title">
                    <I18n id={MessageKey.FrontendTrafficStatsStatusDistribution} />
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
                    theme={this.props.theme}
                />
            </div>
        )
    }

    private renderResponsesTable(zone: ZoneData) {
        const data = [
            { status: "1xx", count: zone.responses["1xx"] },
            { status: "2xx", count: zone.responses["2xx"] },
            { status: "3xx", count: zone.responses["3xx"] },
            { status: "4xx", count: zone.responses["4xx"] },
            { status: "5xx", count: zone.responses["5xx"] },
        ]

        const columns = [
            {
                title: <I18n id={MessageKey.FrontendTrafficStatsResponseStatus} />,
                dataIndex: "status",
                key: "status",
            },
            {
                title: <I18n id={MessageKey.FrontendTrafficStatsConnectionsRequests} />,
                dataIndex: "count",
                key: "count",
                render: (count: number) => formatNumber(count),
            },
        ]

        return (
            <div className="traffic-stats-table-container">
                <p className="traffic-stats-chart-title">
                    <I18n id={MessageKey.FrontendTrafficStatsStatusDistribution} />
                </p>
                <Table dataSource={data} columns={columns} pagination={false} rowKey="status" size="small" />
            </div>
        )
    }

    private renderResponseTimeChart(zone: ZoneData) {
        const data = buildResponseTimeData(zone.requestMsecs)

        return (
            <div className="traffic-stats-chart-container">
                <p className="traffic-stats-chart-title">
                    <I18n id={MessageKey.FrontendTrafficStatsResponseTime} />
                </p>
                {data.length === 0 ? (
                    <Empty description={<I18n id={MessageKey.FrontendTrafficStatsNoData} />} />
                ) : (
                    <Area
                        data={data}
                        xField="time"
                        yField="value"
                        height={300}
                        axis={{ x: { labelAutoHide: true } }}
                        theme={this.props.theme}
                    />
                )}
            </div>
        )
    }

    private renderUserAgentChart() {
        const { filterZones } = this.props.stats
        const { selectedHostId } = this.state
        if (!selectedHostId) return null

        const userAgentZone = filterZones[`userAgent@host:${selectedHostId}`]
        const data = userAgentZone ? buildUserAgentData(userAgentZone) : []

        return (
            <div className="traffic-stats-chart-container">
                <p className="traffic-stats-chart-title">
                    <I18n id={MessageKey.FrontendTrafficStatsUserAgents} />
                </p>
                {data.length === 0 ? (
                    <Empty description={<I18n id={MessageKey.FrontendTrafficStatsNoData} />} />
                ) : (
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
                        theme={this.props.theme}
                    />
                )}
            </div>
        )
    }

    private renderCountryCodeChart() {
        const { filterZones } = this.props.stats
        const { selectedHostId } = this.state
        if (!selectedHostId) return null

        const countryCodeZone = filterZones[`countryCode@host:${selectedHostId}`]
        const data = countryCodeZone ? buildCountryCodeData(countryCodeZone) : []

        return (
            <div className="traffic-stats-chart-container">
                <p className="traffic-stats-chart-title">
                    <I18n id={MessageKey.FrontendTrafficStatsCountryCode} />
                </p>
                {data.length === 0 ? (
                    <Empty description={<I18n id={MessageKey.FrontendTrafficStatsNoData} />} />
                ) : (
                    <Pie
                        data={data}
                        angleField="value"
                        colorField="country"
                        radius={0.8}
                        innerRadius={0.6}
                        label={{
                            text: "country",
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
                )}
            </div>
        )
    }

    render() {
        const zone = this.getSelectedZoneData()

        return (
            <div className="traffic-stats-tab-content">
                {this.renderHostSelector()}
                {zone ? (
                    <>
                        {this.renderStatCards(zone)}
                        <Flex className="traffic-stats-charts-row">
                            {this.renderStatusPieChart(zone)}
                            {this.renderResponsesTable(zone)}
                        </Flex>
                        <Flex className="traffic-stats-charts-row">{this.renderResponseTimeChart(zone)}</Flex>
                        <Flex className="traffic-stats-charts-row">
                            {this.renderUserAgentChart()}
                            {this.renderCountryCodeChart()}
                        </Flex>
                    </>
                ) : (
                    <Empty description={<I18n id={MessageKey.FrontendTrafficStatsSelectHost} />} />
                )}
            </div>
        )
    }
}
