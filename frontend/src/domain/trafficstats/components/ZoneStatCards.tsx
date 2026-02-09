import React from "react"
import { Flex, Statistic } from "antd"
import { I18n } from "../../../core/i18n/I18n"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { formatBytes, formatNumber, formatMs } from "../utils/StatsFormatters"

export interface ZoneStatCardsProps {
    requests: number
    inBytes: number
    outBytes: number
    avgResponseTime: number
}

export default class ZoneStatCards extends React.PureComponent<ZoneStatCardsProps> {
    render() {
        const { requests, inBytes, outBytes, avgResponseTime } = this.props

        return (
            <Flex className="traffic-stats-cards-row">
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficStatsConnectionsRequests} />}
                        value={formatNumber(requests)}
                    />
                </div>
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficStatsBytesReceived} />}
                        value={formatBytes(inBytes)}
                    />
                </div>
                <div className="traffic-stats-stat-card">
                    <Statistic
                        title={<I18n id={MessageKey.FrontendTrafficStatsBytesSent} />}
                        value={formatBytes(outBytes)}
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
}
