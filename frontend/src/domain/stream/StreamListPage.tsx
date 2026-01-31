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
import AccessControl from "../../core/components/accesscontrol/AccessControl"
import { UserAccessLevel } from "../user/model/UserAccessLevel"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import AccessDeniedModal from "../../core/components/accesscontrol/AccessDeniedModal"
import StreamTypeDescription from "./utils/StreamTypeDescription"
import StreamSupportWarning from "./components/StreamSupportWarning"
import MessageKey from "../../core/i18n/model/MessageKey.generated"
import { I18n, i18n, raw } from "../../core/i18n/I18n"

export default class StreamListPage extends React.PureComponent {
    private readonly service: StreamService
    private readonly table: React.RefObject<DataTable<StreamResponse> | null>

    constructor(props: any) {
        super(props)
        this.service = new StreamService()
        this.table = React.createRef()
    }

    private isReadOnlyMode(): boolean {
        return !isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.streams)
    }

    private toggleStreamStatus(stream: StreamResponse) {
        if (this.isReadOnlyMode()) {
            return AccessDeniedModal.show()
        }

        const action = stream.enabled ? i18n(MessageKey.CommonDisable) : i18n(MessageKey.CommonEnable)
        UserConfirmation.ask({ id: MessageKey.FrontendStreamToggleConfirmation, params: { action } })
            .then(() => this.service.toggleEnabled(stream.id))
            .then(() => {
                const msgKey = stream.enabled ? MessageKey.CommonTypeDisabled : MessageKey.CommonTypeEnabled
                Notification.success(
                    { id: msgKey, params: { type: MessageKey.CommonStream } },
                    MessageKey.CommonSuccessMessage,
                )
                ReloadNginxAction.execute()
                this.table.current?.refresh()
            })
            .catch(() => {
                const msgKey = stream.enabled ? MessageKey.CommonUnableToDisable : MessageKey.CommonUnableToEnable
                Notification.error(
                    { id: msgKey, params: { type: MessageKey.CommonStream } },
                    MessageKey.CommonUnexpectedErrorTryAgain,
                )
            })
    }

    private buildColumns(): DataTableColumn<StreamResponse>[] {
        return [
            {
                id: "name",
                description: MessageKey.CommonName,
                renderer: item => item.name,
            },
            {
                id: "binding.type",
                description: MessageKey.CommonType,
                renderer: item => <I18n id={StreamTypeDescription[item.type]} />,
                width: 200,
            },
            {
                id: "binding.address",
                description: MessageKey.CommonBinding,
                renderer: item => `${item.binding.address}:${item.binding.port}`,
                width: 200,
            },
            {
                id: "enabled",
                description: MessageKey.CommonEnabled,
                renderer: item => DataTableRenderers.yesNo(item.enabled),
                width: 150,
            },
            {
                id: "actions",
                description: raw(""),
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
        if (this.isReadOnlyMode()) {
            return AccessDeniedModal.show()
        }

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
            title: MessageKey.CommonStreams,
            subtitle: MessageKey.FrontendStreamListSubtitle,
            actions: [
                {
                    description: MessageKey.FrontendStreamNewButton,
                    onClick: "/streams/new",
                    disabled: this.isReadOnlyMode(),
                },
            ],
        })
    }

    render() {
        return (
            <AccessControl
                requiredAccessLevel={UserAccessLevel.READ_ONLY}
                permissionResolver={permissions => permissions.streams}
            >
                <StreamSupportWarning />

                <DataTable
                    id="streams"
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
