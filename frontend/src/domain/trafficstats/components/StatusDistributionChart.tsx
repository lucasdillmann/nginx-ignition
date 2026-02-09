import React from "react"
import { Empty } from "antd"
import { Pie } from "@ant-design/charts"
import { I18n } from "../../../core/i18n/I18n"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { STATUS_COLORS, StatusDataItem } from "../utils/StatsChartUtils"

export interface StatusDistributionChartProps {
    data: StatusDataItem[]
    theme: "light" | "dark"
    disableAnimation?: boolean
}

export default class StatusDistributionChart extends React.PureComponent<StatusDistributionChartProps> {
    render() {
        const { data, theme, disableAnimation } = this.props

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
                    theme={theme}
                    animation={!disableAnimation}
                />
            </div>
        )
    }
}
