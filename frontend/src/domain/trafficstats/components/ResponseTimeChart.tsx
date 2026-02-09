import React from "react"
import { Empty } from "antd"
import { Area } from "@ant-design/charts"
import { I18n } from "../../../core/i18n/I18n"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"

interface ResponseTimeData {
    time: string
    timestamp: number
    value: number
}

export interface ResponseTimeChartProps {
    data: ResponseTimeData[]
    theme: "light" | "dark"
    disableAnimation?: boolean
}

export default class ResponseTimeChart extends React.PureComponent<ResponseTimeChartProps> {
    render() {
        const { data, theme, disableAnimation } = this.props

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
                        theme={theme}
                        // @ts-expect-error attribute not mapped in the TS contract
                        animation={!disableAnimation}
                    />
                )}
            </div>
        )
    }
}
