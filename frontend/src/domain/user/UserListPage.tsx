import React from "react";
import DataTable, {DataTableColumn} from "../../core/components/datatable/DataTable";
import DataTableRenderers from "../../core/components/datatable/DataTableRenderers";
import {Link} from "react-router-dom";
import {DeleteOutlined, EditOutlined} from "@ant-design/icons";
import PageResponse from "../../core/pagination/PageResponse";
import UserService from "./UserService";
import UserResponse from "./model/UserResponse";
import {UserRole} from "./model/UserRole";
import AppContext from "../../core/components/context/AppContext";
import {Tooltip} from "antd";
import AppShellContext from "../../core/components/shell/AppShellContext";
import DeleteUserAction from "./actions/DeleteUserAction";

export default class UserListPage extends React.PureComponent {
    static contextType = AppShellContext
    context!: React.ContextType<typeof AppShellContext>

    private readonly service: UserService
    private readonly table: React.RefObject<DataTable<UserResponse>>

    constructor(props: any) {
        super(props)
        this.service = new UserService()
        this.table = React.createRef()
    }

    private renderDeleteButton(item: UserResponse, user?: UserResponse): React.ReactNode {
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
                description: "Name",
                renderer: (item) => item.name,
            },
            {
                id: "username",
                description: "Username",
                renderer: (item) => item.username,
                width: 250,
            },
            {
                id: "role",
                description: "Role",
                renderer: (item) => this.translateRole(item.role),
                width: 150,
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
                        <Link to={`/users/${item.id}`}>
                            <EditOutlined className="action-icon" />
                        </Link>

                        <AppContext.Consumer>
                            {context => this.renderDeleteButton(item, context.user)}
                        </AppContext.Consumer>
                    </>
                ),
                fixed: true,
                width: 120,
            }
        ]
    }

    private translateRole(role: UserRole): string {
        switch (role) {
            case UserRole.REGULAR_USER: return "Regular user"
            case UserRole.ADMIN: return "Admin"
        }
    }

    private async deleteUser(user: UserResponse) {
        return DeleteUserAction
            .execute(user.id)
            .then(() => this.table.current?.refresh())
    }

    private fetchData(pageSize: number, pageNumber: number): Promise<PageResponse<UserResponse>> {
        return this.service.list(pageSize, pageNumber)
    }

    componentDidMount() {
        this.context.updateConfig({
            title: "Users",
            subtitle: "Relation of the nginx ignition's users",
            actions: [
                {
                    description: "New user",
                    onClick: "/users/new",
                }
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
