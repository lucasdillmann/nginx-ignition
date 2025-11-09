import React from "react"
import DataTable, { DataTableColumn } from "../../core/components/datatable/DataTable"
import { Flex } from "antd"
import { Link } from "react-router-dom"
import { DeleteOutlined, EditOutlined } from "@ant-design/icons"
import PageResponse from "../../core/pagination/PageResponse"
import Preloader from "../../core/components/preloader/Preloader"
import AppShellContext from "../../core/components/shell/AppShellContext"
import CommonNotifications from "../../core/components/notification/CommonNotifications"
import EmptyStates from "../../core/components/emptystate/EmptyStates"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import { UserAccessLevel } from "../user/model/UserAccessLevel"
import AccessDeniedPage from "../../core/components/accesscontrol/AccessDeniedPage"
import AccessDeniedModal from "../../core/components/accesscontrol/AccessDeniedModal"
import VpnResponse from "./model/VpnResponse"
import VpnService from "./VpnService"
import AvailableDriverResponse from "./model/AvailableDriverResponse"
import DeleteVpnAction from "./actions/DeleteVpnAction"

interface VpnListPageState {
    loading: boolean
    drivers: AvailableDriverResponse[]
    error?: Error
}

export default class VpnListPage extends React.Component<any, VpnListPageState> {
    private readonly service: VpnService
    private readonly table: React.RefObject<DataTable<VpnResponse> | null>

    constructor(props: any) {
        super(props)
        this.service = new VpnService()
        this.table = React.createRef()
        this.state = {
            loading: true,
            drivers: [],
        }
    }

    private isReadOnlyMode() {
        return !isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.vpns)
    }

    private translateDriverName(driverId: string): string {
        const { drivers } = this.state
        return drivers.find(driver => driver.id === driverId)?.name ?? driverId
    }

    private buildColumns(): DataTableColumn<VpnResponse>[] {
        return [
            {
                id: "name",
                description: "Name",
                renderer: item => item.name,
            },
            {
                id: "driver",
                description: "Driver",
                renderer: item => this.translateDriverName(item.driver),
                width: 250,
            },
            {
                id: "actions",
                description: "",
                renderer: item => (
                    <>
                        <Link to={`/vpns/${item.id}`}>
                            <EditOutlined className="action-icon" />
                        </Link>

                        <Link to="" onClick={() => this.deleteVpn(item)}>
                            <DeleteOutlined className="action-icon" />
                        </Link>
                    </>
                ),
                width: 100,
            },
        ]
    }

    private async deleteVpn(vpn: VpnResponse) {
        if (this.isReadOnlyMode()) {
            return AccessDeniedModal.show()
        }

        return DeleteVpnAction.execute(vpn.id).then(() => this.table.current?.refresh())
    }

    private fetchData(pageSize: number, pageNumber: number, searchTerms?: string): Promise<PageResponse<VpnResponse>> {
        return this.service.list(pageSize, pageNumber, searchTerms)
    }

    componentDidMount() {
        this.service
            .availableDrivers()
            .then(drivers =>
                this.setState({
                    loading: false,
                    drivers,
                }),
            )
            .catch(error => {
                CommonNotifications.failedToFetch()
                this.setState({ loading: false, error })
            })

        AppShellContext.get().updateConfig({
            title: "VPN connections",
            subtitle: "Configuration of the nginx ignition VPN connections",
            actions: [
                {
                    description: "New connection",
                    onClick: "/vpns/new",
                    disabled: this.isReadOnlyMode(),
                },
            ],
        })
    }

    render() {
        if (!isAccessGranted(UserAccessLevel.READ_ONLY, permissions => permissions.vpns)) {
            return <AccessDeniedPage />
        }

        const { loading, error } = this.state
        if (loading)
            return (
                <Flex justify="center" align="center">
                    <Preloader loading />
                </Flex>
            )

        if (error !== undefined) return EmptyStates.FailedToFetch

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
