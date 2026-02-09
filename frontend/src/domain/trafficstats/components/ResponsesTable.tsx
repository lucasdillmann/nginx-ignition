import React from "react"
import { Table } from "antd"
import { I18n } from "../../../core/i18n/I18n"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { ZoneData } from "../model/TrafficStatsResponse"
import { formatNumber } from "../utils/StatsFormatters"

export interface ResponsesTableProps {
    responses: ZoneData["responses"]
}

interface StatusRow {
    status: string
    count: number
}

export default class ResponsesTable extends React.PureComponent<ResponsesTableProps> {
    private getColumns() {
        return [
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
    }

    render() {
        const { responses } = this.props

        const data: StatusRow[] = [
            { status: "1xx", count: responses["1xx"] },
            { status: "2xx", count: responses["2xx"] },
            { status: "3xx", count: responses["3xx"] },
            { status: "4xx", count: responses["4xx"] },
            { status: "5xx", count: responses["5xx"] },
        ]

        return (
            <div className="traffic-stats-table-container">
                <p className="traffic-stats-chart-title">
                    <I18n id={MessageKey.FrontendTrafficStatsStatusDistribution} />
                </p>
                <Table dataSource={data} columns={this.getColumns()} pagination={false} rowKey="status" size="small" />
            </div>
        )
    }
}
