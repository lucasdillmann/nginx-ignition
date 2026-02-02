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
import AccessControl from "../../core/components/accesscontrol/AccessControl"
import { UserAccessLevel } from "../user/model/UserAccessLevel"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import MessageKey from "../../core/i18n/model/MessageKey.generated"
import { I18n, raw } from "../../core/i18n/I18n"
import { themedColors } from "../../core/components/theme/ThemedResources"

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
                description: MessageKey.CommonName,
                renderer: item => item.name,
            },
            {
                id: "realm",
                description: MessageKey.FrontendAccesslistRealm,
                renderer: item => item.realm,
                width: 250,
            },
            {
                id: "defaultOutcome",
                description: MessageKey.FrontendAccesslistDefaultOutcome,
                renderer: item => {
                    switch (item.defaultOutcome) {
                        case AccessListOutcome.ALLOW:
                            return <I18n id={MessageKey.FrontendAccesslistOutcomeAllow} />
                        case AccessListOutcome.DENY:
                            return <I18n id={MessageKey.FrontendAccesslistOutcomeDeny} />
                    }
                },
                width: 200,
            },
            {
                id: "satisfyAll",
                description: MessageKey.CommonMode,
                renderer: item => (
                    <I18n
                        id={
                            item.satisfyAll
                                ? MessageKey.FrontendAccesslistModeSatisfyAll
                                : MessageKey.FrontendAccesslistModeSatisfyAny
                        }
                    />
                ),
                width: 150,
            },
            {
                id: "actions",
                description: raw(""),
                renderer: item => (
                    <>
                        <Link to={`/access-lists/${item.id}`}>
                            <EditOutlined className="action-icon" />
                        </Link>

                        <Link to="" onClick={() => this.deleteAccessList(item)}>
                            <DeleteOutlined style={{ color: themedColors().DANGER }} className="action-icon" />
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
            title: MessageKey.CommonAccessLists,
            subtitle: MessageKey.FrontendAccesslistListSubtitle,
            actions: [
                {
                    description: MessageKey.FrontendAccesslistNewButton,
                    onClick: "/access-lists/new",
                    disabled: !isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.accessLists),
                },
            ],
        })
    }

    render() {
        return (
            <AccessControl
                requiredAccessLevel={UserAccessLevel.READ_ONLY}
                permissionResolver={permissions => permissions.accessLists}
            >
                <DataTable
                    id="access-lists"
                    ref={this.table}
                    columns={this.buildColumns()}
                    dataProvider={(pageSize, pageNumber, searchTerms) =>
                        this.fetchData(pageSize, pageNumber, searchTerms)
                    }
                    rowKey={item => item.id}
                />
            </AccessControl>
        )
    }
}
