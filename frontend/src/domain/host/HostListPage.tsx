import React from "react";
import DataTable, {DataTableColumn} from "../../core/components/datatable/DataTable";
import HostResponse from "./model/HostResponse";
import PageResponse from "../../core/pagination/PageResponse";
import HostService from "./HostService";
import DataTableRenderers from "../../core/components/datatable/DataTableRenderers";
import {EditOutlined, PoweroffOutlined, DeleteOutlined} from "@ant-design/icons";
import "./HostListPage.css"
import {Link} from "react-router-dom";
import UserConfirmation from "../../core/components/confirmation/UserConfirmation";
import Notification from "../../core/components/notification/Notification";
import NginxReload from "../../core/components/nginx/NginxReload";
import TagGroup from "../../core/components/taggroup/TagGroup";
import AppShellContext from "../../core/components/shell/AppShellContext";
import DeleteHostAction from "./actions/DeleteHostAction";

export default class HostListPage extends React.PureComponent {
    static contextType = AppShellContext
    context!: React.ContextType<typeof AppShellContext>

    private readonly service: HostService
    private readonly table: React.RefObject<DataTable<HostResponse>>

    constructor(props: any) {
        super(props)
        this.service = new HostService()
        this.table = React.createRef()
    }

    private buildColumns(): DataTableColumn<HostResponse>[] {
        return [
            {
                id: "domainNames",
                description: "Domain names",
                renderer: (item) => <TagGroup values={item.domainNames} />,
            },
            {
                id: "defaultServer",
                description: "Default",
                renderer: (item) => DataTableRenderers.yesNo(item.defaultServer),
                width: 100,
            },
            {
                id: "enabled",
                description: "Enabled",
                renderer: (item) => DataTableRenderers.yesNo(item.enabled),
                width: 100,
            },
            {
                id: "actions",
                description: "",
                renderer: (item) => (
                    <>
                        <Link to={`/hosts/${item.id}`}>
                            <EditOutlined className="action-icon" />
                        </Link>
                        <Link to="" onClick={() => this.toggleHostStatus(item)}>
                            <PoweroffOutlined className="action-icon" />
                        </Link>

                        <Link to="" onClick={() => this.deleteHost(item)}>
                            <DeleteOutlined className="action-icon" />
                        </Link>
                    </>
                ),
                width: 120,
            }
        ]
    }

    private toggleHostStatus(host: HostResponse) {
        const action = host.enabled ? "disable" : "enable"
        UserConfirmation
            .ask(`Do you really want to ${action} the host?`)
            .then(() => this.service.toggleEnabled(host.id))
            .then(() => Notification.success(
                `Host ${action}d`,
                `The host was ${action}d successfully`,
            ))
            .then(() => this.table.current?.refresh())
            .then(() => NginxReload.ask())
            .catch(() => Notification.error(
                `Unable to ${action} the host`,
                `An unexpected error was found while trying to ${action} the host. Please try again later.`,
            ))
    }

    private async deleteHost(host: HostResponse) {
        return DeleteHostAction
            .execute(host.id)
            .then(() => this.table.current?.refresh())
    }

    private fetchData(pageSize: number, pageNumber: number): Promise<PageResponse<HostResponse>> {
        return this.service.list(pageSize, pageNumber)
    }

    componentDidMount() {
        this.context.updateConfig({
            title: "Hosts",
            subtitle: "Relation of all nginx's virtual hosts definitions",
            actions: [
                {
                    description: "New host",
                    onClick: "/hosts/new",
                },
            ],
        })
    }

    render() {
        return (
            <DataTable
                ref={this.table}
                columns={this.buildColumns()}
                dataProvider={(pageSize, pageNumber) => this.fetchData(pageSize, pageNumber)}
                rowKey={item => item.id}
            />
        );
    }
}
