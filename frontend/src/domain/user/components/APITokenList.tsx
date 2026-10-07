import React from "react"
import { Button, Flex, Typography } from "antd"
import { DeleteOutlined, KeyOutlined, PlusOutlined } from "@ant-design/icons"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { I18n, i18n } from "../../../core/i18n/I18n"
import APITokenResponse from "../model/APITokenResponse"
import { formatDateTime } from "../../../core/i18n/I18nDateTime"
import { themedColors } from "../../../core/components/theme/ThemedResources"
import "./APITokenList.css"

export interface APITokenListProps {
    onCreate: () => void
    onRevoke: (token: APITokenResponse) => void
    tokens: APITokenResponse[]
}

export default class APITokenList extends React.Component<APITokenListProps> {
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
                <Button
                    type="primary"
                    icon={<PlusOutlined />}
                    onClick={this.props.onCreate}
                    className="api-tokens-button"
                >
                    <I18n id={MessageKey.FrontendUserTokensNewButton} />
                </Button>
            </Flex>
        )
    }

    private renderExpiration(token: APITokenResponse) {
        if (!token.expiration)
            return <I18n id={MessageKey.FrontendUserTokensNeverExpires} />

        return formatDateTime(token.expiration)
    }

    private renderToken(token: APITokenResponse) {
        return (
            <Flex key={token.id} className="api-tokens-item" align="center">
                <Flex vertical className="api-tokens-item-details">
                    <Typography.Text strong className="api-tokens-item-name">
                        {token.name}
                    </Typography.Text>
                    <Typography.Text type="secondary" className="api-tokens-item-metadata">
                        {i18n(MessageKey.FrontendUserTokensCreatedLabel)} {formatDateTime(token.createdAt)} ·{" "}
                        {i18n(MessageKey.FrontendUserTokensExpirationLabel)} {this.renderExpiration(token)}
                    </Typography.Text>
                </Flex>

                <Button
                    danger
                    type="text"
                    icon={
                        <DeleteOutlined
                            style={{ color: themedColors().DANGER }}
                            className="action-icon"
                            title={i18n(MessageKey.CommonDelete)}
                        />
                    }
                    onClick={() => this.props.onRevoke(token)}
                />
            </Flex>
        )
    }

    render() {
        const { tokens } = this.props

        if (tokens.length === 0) {
            return this.renderEmptyState()
        }

        return (
            <Flex vertical className="api-tokens-container">
                {tokens.map(token => this.renderToken(token))}

                <Flex justify="end" className="api-tokens-footer">
                    <Button type="dashed" icon={<PlusOutlined />} onClick={this.props.onCreate}>
                        <I18n id={MessageKey.FrontendUserTokensNewButton} />
                    </Button>
                </Flex>
            </Flex>
        )
    }
}
