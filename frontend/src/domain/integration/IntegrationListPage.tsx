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
import IntegrationResponse from "./model/IntegrationResponse"
import IntegrationService from "./IntegrationService"
import AvailableDriverResponse from "./model/AvailableDriverResponse"
import DeleteIntegrationAction from "./actions/DeleteIntegrationAction"
import DataTableRenderers from "../../core/components/datatable/DataTableRenderers"
import MessageKey from "../../core/i18n/model/MessageKey.generated"
import { raw } from "../../core/i18n/I18n"

interface IntegrationListPageState {
    loading: boolean
    drivers: AvailableDriverResponse[]
    error?: Error
}

export default class IntegrationListPage extends React.Component<any, IntegrationListPageState> {
    private readonly service: IntegrationService
    private readonly table: React.RefObject<DataTable<IntegrationResponse> | null>

    constructor(props: any) {
        super(props)
        this.service = new IntegrationService()
        this.table = React.createRef()
        this.state = {
            loading: true,
            drivers: [],
        }
    }

    private isReadOnlyMode() {
        return !isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.integrations)
    }

    private translateDriverName(driverId: string): string {
        const { drivers } = this.state
        return drivers.find(driver => driver.id === driverId)?.name ?? driverId
    }

    private buildColumns(): DataTableColumn<IntegrationResponse>[] {
        return [
            {
                id: "name",
                description: MessageKey.CommonName,
                renderer: item => item.name,
            },
            {
                id: "driver",
                description: MessageKey.CommonDriver,
                renderer: item => this.translateDriverName(item.driver),
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
                        <Link to={`/integrations/${item.id}`}>
                            <EditOutlined className="action-icon" />
                        </Link>

                        <Link to="" onClick={() => this.deleteIntegration(item)}>
                            <DeleteOutlined className="action-icon" />
                        </Link>
                    </>
                ),
                width: 100,
            },
        ]
    }

    private async deleteIntegration(integration: IntegrationResponse) {
        if (this.isReadOnlyMode()) {
            return AccessDeniedModal.show()
        }

        return DeleteIntegrationAction.execute(integration.id).then(() => this.table.current?.refresh())
    }

    private fetchData(
        pageSize: number,
        pageNumber: number,
        searchTerms?: string,
    ): Promise<PageResponse<IntegrationResponse>> {
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
            title: MessageKey.CommonIntegrations,
            subtitle: MessageKey.FrontendIntegrationListSubtitle,
            actions: [
                {
                    description: MessageKey.FrontendIntegrationNewButton,
                    onClick: "/integrations/new",
                    disabled: this.isReadOnlyMode(),
                },
            ],
        })
    }

    render() {
        if (!isAccessGranted(UserAccessLevel.READ_ONLY, permissions => permissions.integrations)) {
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
