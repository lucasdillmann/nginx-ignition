import React from "react"
import DataTable, { DataTableColumn } from "../../core/components/datatable/DataTable"
import HostResponse from "./model/HostResponse"
import PageResponse from "../../core/pagination/PageResponse"
import HostService from "./HostService"
import DataTableRenderers from "../../core/components/datatable/DataTableRenderers"
import { CopyOutlined, DeleteOutlined, EditOutlined, PoweroffOutlined } from "@ant-design/icons"
import "./HostListPage.css"
import { Link } from "react-router-dom"
import UserConfirmation from "../../core/components/confirmation/UserConfirmation"
import Notification from "../../core/components/notification/Notification"
import ReloadNginxAction from "../nginx/actions/ReloadNginxAction"
import TagGroup, { TagGroupItem } from "../../core/components/taggroup/TagGroup"
import AppShellContext from "../../core/components/shell/AppShellContext"
import DeleteHostAction from "./actions/DeleteHostAction"
import { UserAccessLevel } from "../user/model/UserAccessLevel"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import AccessControl from "../../core/components/accesscontrol/AccessControl"
import AccessDeniedModal from "../../core/components/accesscontrol/AccessDeniedModal"
import { navigateTo } from "../../core/components/router/AppRouter"
import { Button } from "antd"
import HostSupportWarning from "./components/HostSupportWarning"
import { HostBindingType } from "./model/HostRequest"
import MessageKey from "../../core/i18n/model/MessageKey.generated"
import { i18n, raw } from "../../core/i18n/I18n"

const BUTTON_STYLE = {
    height: "auto",
}

export default class HostListPage extends React.PureComponent {
    private readonly service: HostService
    private readonly table: React.RefObject<DataTable<HostResponse> | null>

    constructor(props: any) {
        super(props)
        this.service = new HostService()
        this.table = React.createRef()
    }

    private isReadOnlyMode(): boolean {
        return !isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.hosts)
    }

    private handleDomainNames(host: HostResponse): string[] | TagGroupItem[] {
        const { domainNames, useGlobalBindings, globalBindings, bindings } = host
        if (!Array.isArray(domainNames) || domainNames.length == 0) return [i18n(MessageKey.CommonDefaultServerLabel)]

        const targetBindings = useGlobalBindings ? globalBindings : bindings
        if (!Array.isArray(targetBindings) || targetBindings.length == 0) return domainNames

        const targetBinding = targetBindings.find(item => item.type == HostBindingType.HTTPS) ?? targetBindings[0]

        return domainNames.map(domainName => ({
            name: domainName,
            url: `${targetBinding.type.toLowerCase()}://${domainName}:${targetBinding.port}`,
        }))
    }

    private buildColumns(): DataTableColumn<HostResponse>[] {
        return [
            {
                id: "domainNames",
                description: MessageKey.CommonDomainNames,
                renderer: item => <TagGroup values={this.handleDomainNames(item)} />,
            },
            {
                id: "defaultServer",
                description: MessageKey.FrontendHostDefault,
                renderer: item => DataTableRenderers.yesNo(item.defaultServer),
                width: 120,
            },
            {
                id: "enabled",
                description: MessageKey.CommonEnabled,
                renderer: item => DataTableRenderers.yesNo(item.enabled),
                width: 120,
            },
            {
                id: "actions",
                description: raw(""),
                renderer: item => (
                    <>
                        <Link to={`/hosts/${item.id}`}>
                            <EditOutlined className="action-icon" />
                        </Link>
                        <Button
                            type="link"
                            onClick={() => this.cloneHost(item)}
                            style={BUTTON_STYLE}
                            icon={<CopyOutlined className="action-icon" />}
                        />
                        <Button
                            type="link"
                            onClick={() => this.toggleHostStatus(item)}
                            style={BUTTON_STYLE}
                            icon={<PoweroffOutlined className="action-icon" />}
                        />
                        <Button
                            type="link"
                            onClick={() => this.deleteHost(item)}
                            style={BUTTON_STYLE}
                            icon={<DeleteOutlined className="action-icon" />}
                        />
                    </>
                ),
                width: 160,
            },
        ]
    }

    private cloneHost(item: HostResponse) {
        if (this.isReadOnlyMode()) {
            return AccessDeniedModal.show()
        }

        navigateTo(`/hosts/new?copyFrom=${item.id}`)
    }

    private toggleHostStatus(host: HostResponse) {
        if (this.isReadOnlyMode()) {
            return AccessDeniedModal.show()
        }

        const action = host.enabled ? i18n(MessageKey.CommonDisable) : i18n(MessageKey.CommonEnable)
        UserConfirmation.ask({ id: MessageKey.FrontendHostToggleConfirmation, params: { action } })
            .then(() => this.service.toggleEnabled(host.id))
            .then(() => {
                const msgKey = host.enabled ? MessageKey.CommonTypeDisabled : MessageKey.CommonTypeEnabled
                Notification.success(
                    { id: msgKey, params: { type: MessageKey.CommonHost } },
                    MessageKey.CommonSuccessMessage,
                )
                ReloadNginxAction.execute()
                this.table.current?.refresh()
            })
            .catch(() => {
                const msgKey = host.enabled ? MessageKey.CommonUnableToDisable : MessageKey.CommonUnableToEnable
                Notification.error(
                    { id: msgKey, params: { type: MessageKey.CommonHost } },
                    MessageKey.CommonUnexpectedErrorTryAgain,
                )
            })
    }

    private async deleteHost(host: HostResponse) {
        if (this.isReadOnlyMode()) {
            return AccessDeniedModal.show()
        }

        return DeleteHostAction.execute(host.id).then(() => this.table.current?.refresh())
    }

    private fetchData(pageSize: number, pageNumber: number, searchTerms?: string): Promise<PageResponse<HostResponse>> {
        return this.service.list(pageSize, pageNumber, searchTerms)
    }

    componentDidMount() {
        AppShellContext.get().updateConfig({
            title: MessageKey.CommonHosts,
            subtitle: MessageKey.FrontendHostListSubtitle,
            actions: [
                {
                    description: MessageKey.FrontendHostNewButton,
                    onClick: "/hosts/new",
                    disabled: this.isReadOnlyMode(),
                },
            ],
        })
    }

    render() {
        return (
            <AccessControl
                requiredAccessLevel={UserAccessLevel.READ_ONLY}
                permissionResolver={permissions => permissions.hosts}
            >
                <HostSupportWarning />

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
