import React from "react"
import { Flex, Empty } from "antd"
import TrafficStatsResponse, { ZoneData } from "../model/TrafficStatsResponse"
import HostService from "../../host/HostService"
import HostResponse from "../../host/model/HostResponse"
import {
    buildStatusDistributionData,
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
import ResponsesTable from "../components/ResponsesTable"
import ZoneStatCards from "../components/ZoneStatCards"
import TagGroup from "../../../core/components/taggroup/TagGroup"
import PaginatedSelect from "../../../core/components/select/PaginatedSelect"
import { ExclamationCircleOutlined } from "@ant-design/icons"

interface ByHostTabProps {
    stats: TrafficStatsResponse
    theme: "light" | "dark"
}

interface ByHostTabState {
    hosts: HostResponse[]
    selectedHost?: HostResponse
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

    private getSelectedZoneData(): ZoneData | undefined {
        const { filterZones } = this.props.stats
        const { selectedHost } = this.state
        if (!selectedHost || !filterZones.hosts) return undefined
        return filterZones.hosts[selectedHost.id]
    }

    private getAvgResponseTime(zone: ZoneData): number {
        if (zone.requestCounter === 0) return 0
        return zone.requestMsecCounter / zone.requestCounter
    }

    private handleDomainNames(domainNames?: string[]): string[] {
        if (Array.isArray(domainNames) && domainNames.length > 0) return domainNames

        return []
    }

    private renderHostSelector() {
        const { selectedHost } = this.state

        return (
            <Flex className="traffic-stats-settings-option">
                <p>
                    <I18n id={MessageKey.CommonHost} />
                </p>
                <PaginatedSelect<HostResponse>
                    placeholder={MessageKey.CommonSelectOne}
                    onChange={host => this.setState({ selectedHost: host })}
                    pageProvider={(pageSize, pageNumber, searchTerms) =>
                        this.hostService.list(pageSize, pageNumber, searchTerms)
                    }
                    value={selectedHost}
                    itemDescription={item =>
                        item.defaultServer ? (
                            <span style={{ fontStyle: "italic", color: "grey" }}>
                                <I18n id={MessageKey.CommonDefaultServerLabel} />
                            </span>
                        ) : (
                            <TagGroup values={this.handleDomainNames(item.domainNames)} />
                        )
                    }
                    itemKey={item => item.id}
                    autoFocus
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
        const { selectedHost } = this.state
        if (!selectedHost) return null

        const userAgentZone = filterZones[`userAgent@host:${selectedHost.id}`]
        const data = userAgentZone ? buildUserAgentData(userAgentZone) : []
        return <UserAgentChart data={data} theme={this.props.theme} />
    }

    private renderCountryCodeChart() {
        const { filterZones } = this.props.stats
        const { selectedHost } = this.state
        if (!selectedHost) return null

        const countryCodeZone = filterZones[`countryCode@host:${selectedHost.id}`]
        const data = countryCodeZone ? buildCountryCodeData(countryCodeZone) : []
        return <CountryCodeChart data={data} theme={this.props.theme} />
    }

    private renderCityChart() {
        const { filterZones } = this.props.stats
        const { selectedHost } = this.state
        if (!selectedHost) return null

        const cityZone = filterZones[`city@host:${selectedHost.id}`]
        const data = cityZone ? buildCityData(cityZone) : []
        return <CityChart data={data} theme={this.props.theme} />
    }

    private renderMainContents() {
        const { selectedHost } = this.state
        if (!selectedHost) {
            const icon = (
                <ExclamationCircleOutlined style={{ fontSize: 70, color: "var(--nginxIgnition-colorTextDisabled)" }} />
            )

            return <Empty image={icon} description={<I18n id={MessageKey.FrontendTrafficStatsSelectHost} />} />
        }

        const zone = this.getSelectedZoneData()
        if (!zone) return <Empty description={<I18n id={MessageKey.FrontendTrafficStatsNoData} />} />

        return (
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
                    {this.renderCityChart()}
                </Flex>
            </>
        )
    }

    render() {
        return (
            <div className="traffic-stats-tab-content">
                {this.renderHostSelector()}
                {this.renderMainContents()}
            </div>
        )
    }
}
