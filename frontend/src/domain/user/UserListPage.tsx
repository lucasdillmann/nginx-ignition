import React from "react";
import DataTable, {DataTableColumn} from "../../core/components/datatable/DataTable";
import DataTableRenderers from "../../core/components/datatable/DataTableRenderers";
import {Link} from "react-router-dom";
import {DeleteOutlined, EditOutlined} from "@ant-design/icons";
import UserConfirmation from "../../core/components/confirmation/UserConfirmation";
import Notification from "../../core/components/notification/Notification";
import PageResponse from "../../core/pagination/PageResponse";
import UserService from "./UserService";
import UserResponse from "./model/UserResponse";
import {UserRole} from "./model/UserRole";
import AppContext from "../../core/components/context/AppContext";
import {Tooltip} from "antd";
import ShellAwareComponent, {ShellConfig} from "../../core/components/shell/ShellAwareComponent";

export default class UserListPage extends ShellAwareComponent {
    static contextType = AppContext
    context!: React.ContextType<typeof AppContext>
    private readonly service: UserService
    private readonly table: React.RefObject<DataTable<UserResponse>>

    constructor(props: any) {
        super(props)
        this.service = new UserService()
        this.table = React.createRef()
    }

    private renderDeleteButton(item: UserResponse): React.ReactNode {
        const {user} = this.context
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

                        {this.renderDeleteButton(item)}
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

    private deleteUser(user: UserResponse) {
        UserConfirmation
            .ask("Do you really want to delete the user?")
            .then(() => this.service.delete(user.id))
            .then(() => Notification.success(
                `User deleted`,
                `The user was deleted successfully`,
            ))
            .then(() => this.table.current?.refresh())
            .catch(() => Notification.error(
                `Unable to delete the user`,
                `An unexpected error was found while trying to delete the user. Please try again later.`,
            ))
    }

    private fetchData(pageSize: number, pageNumber: number): Promise<PageResponse<UserResponse>> {
        return this.service.list(pageSize, pageNumber)
    }

    shellConfig(): ShellConfig {
        return {
            title: "Users",
            subtitle: "Relation of the nginx ignition's users",
            actions: [
                {
                    description: "New user",
                    onClick: "/users/new",
                }
            ],
        };
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
