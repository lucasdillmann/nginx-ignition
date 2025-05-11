import React from "react"
import DataTable, { DataTableColumn } from "../../core/components/datatable/DataTable"
import { Link } from "react-router-dom"
import { DeleteOutlined, EditOutlined, PoweroffOutlined } from "@ant-design/icons"
import PageResponse from "../../core/pagination/PageResponse"
import StreamService from "./StreamService"
import AppShellContext from "../../core/components/shell/AppShellContext"
import StreamResponse from "./model/StreamResponse"
import DataTableRenderers from "../../core/components/datatable/DataTableRenderers"
import DeleteStreamAction from "./actions/DeleteStreamAction"
import UserConfirmation from "../../core/components/confirmation/UserConfirmation"
import Notification from "../../core/components/notification/Notification"
import ReloadNginxAction from "../nginx/actions/ReloadNginxAction"

export default class StreamListPage extends React.PureComponent {
    private readonly service: StreamService
    private readonly table: React.RefObject<DataTable<StreamResponse> | null>

    constructor(props: any) {
        super(props)
        this.service = new StreamService()
        this.table = React.createRef()
    }

    private toggleStreamStatus(stream: StreamResponse) {
        const action = stream.enabled ? "disable" : "enable"
        UserConfirmation.ask(`Do you really want to ${action} the stream?`)
            .then(() => this.service.toggleEnabled(stream.id))
            .then(() => {
                Notification.success(`Host ${action}d`, `The stream was ${action}d successfully`)
                ReloadNginxAction.execute()
                this.table.current?.refresh()
            })
            .catch(() =>
                Notification.error(
                    `Unable to ${action} the stream`,
                    `An unexpected error was found while trying to ${action} the stream. Please try again later.`,
                ),
            )
    }

    private buildColumns(): DataTableColumn<StreamResponse>[] {
        return [
            {
                id: "name",
                description: "Name",
                renderer: item => item.name,
            },
            {
                id: "binding.address",
                description: "Binding",
                renderer: item => `${item.binding.address}:${item.binding.port}`,
                width: 250,
            },
            {
                id: "binding.protocol",
                description: "Protocol",
                renderer: item => item.binding.protocol,
                width: 150,
            },
            {
                id: "enabled",
                description: "Enabled",
                renderer: item => DataTableRenderers.yesNo(item.enabled),
                width: 150,
            },
            {
                id: "actions",
                description: "",
                renderer: item => (
                    <>
                        <Link to={`/streams/${item.id}`}>
                            <EditOutlined className="action-icon" />
                        </Link>

                        <Link to="" onClick={() => this.toggleStreamStatus(item)}>
                            <PoweroffOutlined className="action-icon" />
                        </Link>

                        <Link to="" onClick={() => this.deleteStream(item)}>
                            <DeleteOutlined className="action-icon" />
                        </Link>
                    </>
                ),
                width: 120,
            },
        ]
    }

    private async deleteStream(stream: StreamResponse) {
        return DeleteStreamAction.execute(stream.id).then(() => this.table.current?.refresh())
    }

    private fetchData(
        pageSize: number,
        pageNumber: number,
        searchTerms?: string,
    ): Promise<PageResponse<StreamResponse>> {
        return this.service.list(pageSize, pageNumber, searchTerms)
    }

    componentDidMount() {
        AppShellContext.get().updateConfig({
            title: "Streams",
            subtitle: "Relation of nginx's raw TCP, UDP and Unix socket proxies",
            actions: [
                {
                    description: "New stream",
                    onClick: "/streams/new",
                },
            ],
        })
    }

    render() {
        return (
            <DataTable
                ref={this.table}
                columns={this.buildColumns()}
                dataProvider={(pageSize, pageNumber, searchTerms) => this.fetchData(pageSize, pageNumber, searchTerms)}
                rowKey={item => item.id}
            />
        )
    }
}
