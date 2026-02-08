import React from "react"
import { Flex, Select, Empty } from "antd"
import TrafficStatsResponse, { ZoneData } from "../model/TrafficStatsResponse"
import HostService from "../../host/HostService"
import HostResponse from "../../host/model/HostResponse"
import {
    buildStatusDistributionData,
    buildResponseTimeData,
    buildUserAgentData,
    buildCountryCodeData,
} from "../utils/StatsChartUtils"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { I18n, i18n } from "../../../core/i18n/I18n"
import UserAgentChart from "../components/UserAgentChart"
import CountryCodeChart from "../components/CountryCodeChart"
import ResponseTimeChart from "../components/ResponseTimeChart"
import StatusDistributionChart from "../components/StatusDistributionChart"
import ResponsesTable from "../components/ResponsesTable"
import ZoneStatCards from "../components/ZoneStatCards"

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
            <ZoneStatCards
                requests={zone.requestCounter}
                inBytes={zone.inBytes}
                outBytes={zone.outBytes}
                avgResponseTime={this.getAvgResponseTime(zone)}
            />
        )
    }

    private renderStatusPieChart(zone: ZoneData) {
        const data = buildStatusDistributionData(zone.responses)
        return <StatusDistributionChart data={data} theme={this.props.theme} />
    }

    private renderResponsesTable(zone: ZoneData) {
        return <ResponsesTable responses={zone.responses} />
    }

    private renderResponseTimeChart(zone: ZoneData) {
        const data = buildResponseTimeData(zone.requestMsecs)
        return <ResponseTimeChart data={data} theme={this.props.theme} />
    }

    private renderUserAgentChart() {
        const { filterZones } = this.props.stats
        const { selectedHostId } = this.state
        if (!selectedHostId) return null

        const userAgentZone = filterZones[`userAgent@host:${selectedHostId}`]
        const data = userAgentZone ? buildUserAgentData(userAgentZone) : []
        return <UserAgentChart data={data} theme={this.props.theme} />
    }

    private renderCountryCodeChart() {
        const { filterZones } = this.props.stats
        const { selectedHostId } = this.state
        if (!selectedHostId) return null

        const countryCodeZone = filterZones[`countryCode@host:${selectedHostId}`]
        const data = countryCodeZone ? buildCountryCodeData(countryCodeZone) : []
        return <CountryCodeChart data={data} theme={this.props.theme} />
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
