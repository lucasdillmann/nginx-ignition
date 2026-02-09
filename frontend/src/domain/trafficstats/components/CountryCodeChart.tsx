import React from "react"
import { Empty } from "antd"
import { Pie } from "@ant-design/charts"
import { I18n } from "../../../core/i18n/I18n"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"

interface CountryCodeData {
    country: string
    value: number
}

export interface CountryCodeChartProps {
    data: CountryCodeData[]
    theme: "light" | "dark"
    disableAnimation?: boolean
}

export default class CountryCodeChart extends React.PureComponent<CountryCodeChartProps> {
    render() {
        const { data, theme, disableAnimation } = this.props

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
                        theme={theme}
                        animation={!disableAnimation}
                    />
                )}
            </div>
        )
    }
}
