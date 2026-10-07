import React from "react"
import { Button, Flex, Typography } from "antd"
import { DeleteOutlined, KeyOutlined, PlusOutlined } from "@ant-design/icons"
import MessageKey from "../../../../../../core/i18n/model/MessageKey.generated"
import { I18n, i18n, raw } from "../../../../../../core/i18n/I18n"
import { formatDateTime } from "../../../../../../core/i18n/I18nDateTime"
import PageResponse from "../../../../../../core/pagination/PageResponse"
import DataTable, { DataTableColumn } from "../../../../../../core/components/datatable/DataTable"
import { themedColors } from "../../../../../../core/components/theme/ThemedResources"
import DeleteAPITokenAction from "../../../../actions/DeleteAPITokenAction"
import UserService from "../../../../UserService"
import APITokenResponse from "../../../../model/APITokenResponse"
import APITokenCreateModal from "./APITokenCreateModal"
import "./UserAPITokensTab.css"

interface UserAPITokensTabState {
    createModalOpen: boolean
}

export default class UserAPITokensTab extends React.Component<Record<string, never>, UserAPITokensTabState> {
    private readonly service: UserService
    private readonly table: React.RefObject<DataTable<APITokenResponse> | null>

    constructor(props: Record<string, never>) {
        super(props)
        this.service = new UserService()
        this.table = React.createRef()
        this.state = { createModalOpen: false }
    }

    private async fetchData(
        pageSize: number,
        pageNumber: number,
        searchTerms?: string,
    ): Promise<PageResponse<APITokenResponse>> {
        return this.service.listTokens(pageSize, pageNumber, searchTerms)
    }

    private buildColumns(): DataTableColumn<APITokenResponse>[] {
        return [
            {
                id: "name",
                description: MessageKey.CommonName,
                renderer: item => item.name,
            },
            {
                id: "expiration",
                description: MessageKey.FrontendUserTokensExpirationLabel,
                renderer: item =>
                    item.expiration ? (
                        formatDateTime(item.expiration)
                    ) : (
                        <I18n id={MessageKey.FrontendUserTokensNeverExpires} />
                    ),
                width: 220,
            },
            {
                id: "actions",
                description: raw(""),
                renderer: item => (
                    <DeleteOutlined
                        style={{ color: themedColors().DANGER }}
                        className="action-icon"
                        title={i18n(MessageKey.CommonDelete)}
                        onClick={() => this.revokeToken(item)}
                    />
                ),
                width: 80,
            },
        ]
    }

    private revokeToken(token: APITokenResponse) {
        DeleteAPITokenAction.execute(token.id).then(() => this.table.current?.refresh())
    }

    private renderEmptyState() {
        return (
            <Flex vertical align="center" justify="center" className="api-tokens-empty-container">
                <KeyOutlined className="api-tokens-empty-icon" />
                <Typography.Title level={4} className="api-tokens-empty-title">
                    <I18n id={MessageKey.FrontendUserTokensEmptyTitle} />
                </Typography.Title>
                <Typography.Text type="secondary" className="api-tokens-empty-description">
                    <I18n id={MessageKey.FrontendUserTokensEmptyDescription} />
                </Typography.Text>
            </Flex>
        )
    }

    render() {
        const { createModalOpen } = this.state

        return (
            <>
                <DataTable<APITokenResponse>
                    id="user-api-tokens"
                    ref={this.table}
                    columns={this.buildColumns()}
                    dataProvider={this.fetchData.bind(this)}
                    rowKey={item => item.id}
                    emptyState={this.renderEmptyState()}
                />

                <Flex justify="end" className="api-tokens-footer">
                    <Button
                        type="primary"
                        icon={<PlusOutlined />}
                        onClick={() => this.setState({ createModalOpen: true })}
                    >
                        <I18n id={MessageKey.FrontendUserTokensNewButton} />
                    </Button>
                </Flex>

                <APITokenCreateModal
                    open={createModalOpen}
                    onCancel={() => this.setState({ createModalOpen: false })}
                    onCreated={() => this.table.current?.refresh()}
                />
            </>
        )
    }
}
