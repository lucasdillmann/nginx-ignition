import React from "react"
import { Flex, Select, Empty } from "antd"
import TrafficStatsResponse, { ZoneData } from "../model/TrafficStatsResponse"
import {
    buildStatusDistributionData,
    buildResponseTimeData,
    buildUserAgentData,
    buildCountryCodeData,
} from "../utils/StatsChartUtils"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { I18n } from "../../../core/i18n/I18n"
import UserAgentChart from "../components/UserAgentChart"
import CountryCodeChart from "../components/CountryCodeChart"
import ResponseTimeChart from "../components/ResponseTimeChart"
import StatusDistributionChart from "../components/StatusDistributionChart"
import ResponsesTable from "../components/ResponsesTable"
import ZoneStatCards from "../components/ZoneStatCards"

interface ByDomainTabProps {
    stats: TrafficStatsResponse
    theme: "light" | "dark"
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
        const { selectedDomain } = this.state
        if (!selectedDomain) return null

        const userAgentZone = filterZones[`userAgent@domain:${selectedDomain}`]
        const data = userAgentZone ? buildUserAgentData(userAgentZone) : []
        return <UserAgentChart data={data} theme={this.props.theme} />
    }

    private renderCountryCodeChart() {
        const { filterZones } = this.props.stats
        const { selectedDomain } = this.state
        if (!selectedDomain) return null

        const countryCodeZone = filterZones[`countryCode@domain:${selectedDomain}`]
        const data = countryCodeZone ? buildCountryCodeData(countryCodeZone) : []
        return <CountryCodeChart data={data} theme={this.props.theme} />
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
                        <Flex className="traffic-stats-charts-row">
                            {this.renderUserAgentChart()}
                            {this.renderCountryCodeChart()}
                        </Flex>
                    </>
                ) : (
                    <Empty description={<I18n id={MessageKey.FrontendTrafficStatsSelectDomain} />} />
                )}
            </div>
        )
    }
}
