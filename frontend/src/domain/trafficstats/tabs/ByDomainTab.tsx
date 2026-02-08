import React from "react"
import { Flex, Select, Statistic, Empty, Table } from "antd"
import { Pie, Area } from "@ant-design/charts"
import TrafficStatsResponse, { ZoneData } from "../model/TrafficStatsResponse"
import { formatBytes, formatNumber, formatMs } from "../utils/StatsFormatters"
import { buildStatusDistributionData, STATUS_COLORS, buildResponseTimeData } from "../utils/StatsChartUtils"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { I18n } from "../../../core/i18n/I18n"

interface ByDomainTabProps {
    stats: TrafficStatsResponse
}

interface ByDomainTabState {
    selectedDomain?: string
}

export default class ByDomainTab extends React.Component<ByDomainTabProps, ByDomainTabState> {
    constructor(props: ByDomainTabProps) {
        super(props)
        this.state = {}
    }

    private getSelectedZoneData(): ZoneData | undefined {
        const { serverZones } = this.props.stats
        const { selectedDomain } = this.state
        if (!selectedDomain) return undefined
        return serverZones[selectedDomain]
    }

    private getAvgResponseTime(zone: ZoneData): number {
        if (zone.requestCounter === 0) return 0
        return zone.requestMsecCounter / zone.requestCounter
    }

    private renderDomainSelector() {
        const { serverZones } = this.props.stats
        const { selectedDomain } = this.state
        const domains = Object.keys(serverZones).filter(d => d !== "*")

        const options = domains.map(domain => ({
            value: domain,
            label: domain,
        }))

        return (
            <Flex className="traffic-stats-settings-option">
                <p>
                    <I18n id={MessageKey.CommonDomain} />
                </p>
                <Select
                    className="traffic-stats-selector"
                    placeholder={<I18n id={MessageKey.FrontendTrafficStatsSelectDomain} />}
                    options={options}
                    value={selectedDomain}
                    onChange={value => this.setState({ selectedDomain: value })}
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

        if (data.length === 0) {
            return null
        }

        return (
            <div className="traffic-stats-chart-container">
                <p className="traffic-stats-chart-title">
                    <I18n id={MessageKey.FrontendTrafficStatsResponseTime} />
                </p>
                <Area data={data} xField="time" yField="value" height={300} axis={{ x: { labelAutoHide: true } }} />
            </div>
        )
    }

    render() {
        const zone = this.getSelectedZoneData()

        return (
            <div className="traffic-stats-tab-content">
                {this.renderDomainSelector()}
                {zone ? (
                    <>
                        {this.renderStatCards(zone)}
                        <Flex className="traffic-stats-charts-row">
                            {this.renderStatusPieChart(zone)}
                            {this.renderResponsesTable(zone)}
                        </Flex>
                        <Flex className="traffic-stats-charts-row">{this.renderResponseTimeChart(zone)}</Flex>
                    </>
                ) : (
                    <Empty description={<I18n id={MessageKey.FrontendTrafficStatsSelectDomain} />} />
                )}
            </div>
        )
    }
}
