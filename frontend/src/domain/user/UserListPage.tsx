import React from "react"
import DataTable, { DataTableColumn } from "../../core/components/datatable/DataTable"
import DataTableRenderers from "../../core/components/datatable/DataTableRenderers"
import { Link } from "react-router-dom"
import { DeleteOutlined, EditOutlined } from "@ant-design/icons"
import PageResponse from "../../core/pagination/PageResponse"
import UserService from "./UserService"
import UserResponse from "./model/UserResponse"
import AppContext from "../../core/components/context/AppContext"
import { Tooltip } from "antd"
import AppShellContext from "../../core/components/shell/AppShellContext"
import DeleteUserAction from "./actions/DeleteUserAction"
import { UserAccessLevel } from "./model/UserAccessLevel"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import AccessControl from "../../core/components/accesscontrol/AccessControl"
import AccessDeniedModal from "../../core/components/accesscontrol/AccessDeniedModal"
import MessageKey from "../../core/i18n/model/MessageKey.generated"
import { raw } from "../../core/i18n/I18n"

export default class UserListPage extends React.PureComponent {
    private readonly service: UserService
    private readonly table: React.RefObject<DataTable<UserResponse> | null>

    constructor(props: any) {
        super(props)
        this.service = new UserService()
        this.table = React.createRef()
    }

    private isReadOnlyMode(): boolean {
        return !isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.users)
    }

    private renderDeleteButton(item: UserResponse): React.ReactNode {
        const { user } = AppContext.get()
        if (user?.id !== item.id)
            return (
                <Link to="" onClick={() => this.deleteUser(item)}>
                    <DeleteOutlined className="action-icon" />
                </Link>
            )

        return (
            <Tooltip title="You can't delete your own user">
                <DeleteOutlined className="action-icon" disabled />
            </Tooltip>
        )
    }

    private buildColumns(): DataTableColumn<UserResponse>[] {
        return [
            {
                id: "name",
                description: MessageKey.CommonName,
                renderer: item => item.name,
            },
            {
                id: "username",
                description: MessageKey.CommonUsername,
                renderer: item => item.username,
                width: 250,
            },
            {
                id: "enabled",
                description: MessageKey.CommonEnabled,
                renderer: item => DataTableRenderers.yesNo(item.enabled),
                width: 100,
            },
            {
                id: "actions",
                description: raw(""),
                renderer: item => (
                    <>
                        <Link to={`/users/${item.id}`}>
                            <EditOutlined className="action-icon" />
                        </Link>

                        {this.renderDeleteButton(item)}
                    </>
                ),
                width: 120,
            },
        ]
    }

    private async deleteUser(user: UserResponse) {
        if (this.isReadOnlyMode()) {
            return AccessDeniedModal.show()
        }

        return DeleteUserAction.execute(user.id).then(() => this.table.current?.refresh())
    }

    private fetchData(pageSize: number, pageNumber: number, searchTerms?: string): Promise<PageResponse<UserResponse>> {
        return this.service.list(pageSize, pageNumber, searchTerms)
    }

    componentDidMount() {
        AppShellContext.get().updateConfig({
            title: MessageKey.CommonUsers,
            subtitle: MessageKey.FrontendUserListSubtitle,
            actions: [
                {
                    description: MessageKey.FrontendUserNewButton,
                    onClick: "/users/new",
                    disabled: this.isReadOnlyMode(),
                },
            ],
        })
    }

    render() {
        return (
            <AccessControl
                requiredAccessLevel={UserAccessLevel.READ_ONLY}
                permissionResolver={permissions => permissions.users}
            >
                <DataTable
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
