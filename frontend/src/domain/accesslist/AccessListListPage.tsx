import React from "react"
import DataTable, { DataTableColumn } from "../../core/components/datatable/DataTable"
import { Link } from "react-router-dom"
import { DeleteOutlined, EditOutlined } from "@ant-design/icons"
import PageResponse from "../../core/pagination/PageResponse"
import AccessListService from "./AccessListService"
import AppShellContext from "../../core/components/shell/AppShellContext"
import AccessListResponse from "./model/AccessListResponse"
import DeleteAccessListAction from "./actions/DeleteAccessListAction"
import { AccessListOutcome } from "./model/AccessListRequest"

export default class AccessListListPage extends React.PureComponent {
    private readonly service: AccessListService
    private readonly table: React.RefObject<DataTable<AccessListResponse> | null>

    constructor(props: any) {
        super(props)
        this.service = new AccessListService()
        this.table = React.createRef()
    }

    private buildColumns(): DataTableColumn<AccessListResponse>[] {
        return [
            {
                id: "name",
                description: "Name",
                renderer: item => item.name,
            },
            {
                id: "realm",
                description: "Realm",
                renderer: item => item.realm,
                width: 250,
            },
            {
                id: "defaultOutcome",
                description: "Default outcome",
                renderer: item => {
                    switch (item.defaultOutcome) {
                        case AccessListOutcome.ALLOW:
                            return "Allow access"
                        case AccessListOutcome.DENY:
                            return "Deny access"
                    }
                },
                width: 150,
            },
            {
                id: "satisfyAll",
                description: "Mode",
                renderer: item => (item.satisfyAll ? "Satisfy all" : "Satisfy any"),
                width: 150,
            },
            {
                id: "actions",
                description: "",
                renderer: item => (
                    <>
                        <Link to={`/access-lists/${item.id}`}>
                            <EditOutlined className="action-icon" />
                        </Link>

                        <Link to="" onClick={() => this.deleteAccessList(item)}>
                            <DeleteOutlined className="action-icon" />
                        </Link>
                    </>
                ),
                width: 120,
            },
        ]
    }

    private async deleteAccessList(accessList: AccessListResponse) {
        return DeleteAccessListAction.execute(accessList.id).then(() => this.table.current?.refresh())
    }

    private fetchData(
        pageSize: number,
        pageNumber: number,
        searchTerms?: string,
    ): Promise<PageResponse<AccessListResponse>> {
        return this.service.list(pageSize, pageNumber, searchTerms)
    }

    componentDidMount() {
        AppShellContext.get().updateConfig({
            title: "Access lists",
            subtitle: "Relation of the access lists for the nginx authentication and access control",
            actions: [
                {
                    description: "New access list",
                    onClick: "/access-lists/new",
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
