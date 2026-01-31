import React from "react"
import DataTable, { DataTableColumn } from "../../core/components/datatable/DataTable"
import { Link } from "react-router-dom"
import { DeleteOutlined, EditOutlined } from "@ant-design/icons"
import PageResponse from "../../core/pagination/PageResponse"
import CacheService from "./CacheService"
import AppShellContext from "../../core/components/shell/AppShellContext"
import CacheResponse from "./model/CacheResponse"
import DeleteCacheAction from "./actions/DeleteCacheAction"
import AccessControl from "../../core/components/accesscontrol/AccessControl"
import { UserAccessLevel } from "../user/model/UserAccessLevel"
import { isAccessGranted } from "../../core/components/accesscontrol/IsAccessGranted"
import TagGroup from "../../core/components/taggroup/TagGroup"
import MessageKey from "../../core/i18n/model/MessageKey.generated"
import { I18n, raw } from "../../core/i18n/I18n"
import { themedColors } from "../../core/components/theme/ThemedResources"

export default class CacheListPage extends React.PureComponent {
    private readonly service: CacheService
    private readonly table: React.RefObject<DataTable<CacheResponse> | null>

    constructor(props: any) {
        super(props)
        this.service = new CacheService()
        this.table = React.createRef()
    }

    private buildColumns(): DataTableColumn<CacheResponse>[] {
        return [
            {
                id: "name",
                description: MessageKey.CommonName,
                renderer: item => item.name,
            },
            {
                id: "fileExtensions",
                description: MessageKey.CommonFileExtensions,
                renderer: item =>
                    item.fileExtensions ? (
                        <TagGroup values={item.fileExtensions} maximumSize={3} />
                    ) : (
                        <I18n id={MessageKey.FrontendCacheListAllExtensions} />
                    ),
                width: 300,
            },
            {
                id: "maximumSizeMb",
                description: MessageKey.CommonMaximumSize,
                renderer: item =>
                    item.maximumSizeMb ? (
                        <>
                            {item.maximumSizeMb} <I18n id={MessageKey.CommonUnitMb} />
                        </>
                    ) : (
                        <I18n id={MessageKey.FrontendCacheListUnlimited} />
                    ),
                width: 150,
            },
            {
                id: "actions",
                description: raw(""),
                renderer: item => (
                    <>
                        <Link to={`/caches/${item.id}`}>
                            <EditOutlined className="action-icon" />
                        </Link>

                        <Link to="" onClick={() => this.deleteCache(item)}>
                            <DeleteOutlined style={{ color: themedColors().DANGER }} className="action-icon" />
                        </Link>
                    </>
                ),
                width: 120,
            },
        ]
    }

    private async deleteCache(cache: CacheResponse) {
        return DeleteCacheAction.execute(cache.id).then(() => this.table.current?.refresh())
    }

    private fetchData(
        pageSize: number,
        pageNumber: number,
        searchTerms?: string,
    ): Promise<PageResponse<CacheResponse>> {
        return this.service.list(pageSize, pageNumber, searchTerms)
    }

    componentDidMount() {
        AppShellContext.get().updateConfig({
            title: MessageKey.CommonCacheConfigurations,
            subtitle: MessageKey.FrontendCacheListSubtitle,
            actions: [
                {
                    description: MessageKey.FrontendCacheNewButton,
                    onClick: "/caches/new",
                    disabled: !isAccessGranted(UserAccessLevel.READ_WRITE, permissions => permissions.caches),
                },
            ],
        })
    }

    render() {
        return (
            <AccessControl
                requiredAccessLevel={UserAccessLevel.READ_ONLY}
                permissionResolver={permissions => permissions.caches}
            >
                <DataTable
                    id="cache-configurations"
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
