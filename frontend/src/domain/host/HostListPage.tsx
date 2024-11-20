import React from "react";
import DataTable, {DataTableColumn} from "../../core/components/datatable/DataTable";
import HostResponse from "./model/HostResponse";
import PageResponse from "../../core/pagination/PageResponse";
import HostService from "./HostService";
import DataTableRenderers from "../../core/components/datatable/DataTableRenderers";
import {Tag} from "antd";
import {EditOutlined, PoweroffOutlined, DeleteOutlined} from "@ant-design/icons";
import "./HostListPage.css"
import {Link} from "react-router-dom";
import UserConfirmation from "../../core/components/confirmation/UserConfirmation";
import NotificationFacade from "../../core/components/notification/NotificationFacade";
import NginxReload from "../../core/components/nginx/NginxReload";

const MAXIMUM_DOMAIN_NAMES = 3

export default class HostListPage extends React.PureComponent {
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
                renderer: ({domainNames}) => {
                    const tags = domainNames
                        .slice(0, MAXIMUM_DOMAIN_NAMES)
                        .map(name => <Tag>{name}</Tag>)
                    const additionalMessage = domainNames.length > MAXIMUM_DOMAIN_NAMES
                        ? `and ${domainNames.length - MAXIMUM_DOMAIN_NAMES} more`
                        : ""

                    return (
                        <>{tags} {additionalMessage}</>
                    )
                },
            },
            {
                id: "default",
                description: "Default",
                renderer: (item) => DataTableRenderers.yesNo(item.default),
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
                fixed: true,
                width: 120,
            }
        ]
    }

    private toggleHostStatus(host: HostResponse) {
        const action = host.enabled ? "disable" : "enable"
        UserConfirmation
            .ask(`Do you really want to ${action} the host?`)
            .then(() => this.service.toggleEnabled(host.id))
            .then(() => NotificationFacade.success(
                `Host ${action}d`,
                `The host was ${action}d successfully`,
            ))
            .then(() => this.table.current?.refresh())
            .then(() => NginxReload.ask())
            .catch(() => NotificationFacade.error(
                `Unable to ${action} the host`,
                `An unexpected error was found while trying to ${action} the host. Please try again later.`,
            ))
    }

    private deleteHost(host: HostResponse) {
        UserConfirmation
            .ask("Do you really want to delete the host?")
            .then(() => this.service.delete(host.id))
            .then(() => NotificationFacade.success(
                `Host deleted`,
                `The host was deleted successfully`,
            ))
            .then(() => this.table.current?.refresh())
            .then(() => NginxReload.ask())
            .catch(() => NotificationFacade.error(
                `Unable to delete the host`,
                `An unexpected error was found while trying to delete the host. Please try again later.`,
            ))
    }

    private fetchData(pageSize: number, pageNumber: number): Promise<PageResponse<HostResponse>> {
        return this.service.list(pageSize, pageNumber)
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
