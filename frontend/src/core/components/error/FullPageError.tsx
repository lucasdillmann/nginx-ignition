import React from "react"
import { Button, Flex } from "antd"
import "./FullPageError.css"
import { ExclamationCircleFilled } from "@ant-design/icons"
import { themedModal } from "../theme/ThemedResources"
import MessageKey from "../../i18n/model/MessageKey.generated"
import { I18n, i18n } from "../../i18n/I18n"

export interface FullPageErrorProps {
    title?: MessageKey
    message?: MessageKey
    error?: Error
    icon?: any
}

export default class FullPageError extends React.Component<FullPageErrorProps> {
    private openErrorDetailsModal() {
        const { error } = this.props
        themedModal().info({
            title: i18n(MessageKey.FrontendComponentsErrorDetails, {}, "Error details"),
            type: "info",
            width: 1000,
            content: <pre>{error?.stack ?? error?.message ?? typeof error}</pre>,
        })
    }

    private renderErrorButton() {
        const { error } = this.props
        if (error === undefined) return null

        return (
            <Button
                variant="outlined"
                onClick={() => this.openErrorDetailsModal()}
                size="small"
                className="error-page-details-button"
            >
                <I18n id={MessageKey.FrontendComponentsErrorShowDetailsButton} fallback="Show details of what happened" />
            </Button>
        )
    }

    private renderIcon() {
        const { icon } = this.props
        if (icon !== undefined) return icon

        return <ExclamationCircleFilled style={{ fontSize: 48, color: "var(--nginxIgnition-colorError)" }} />
    }

    render() {
        const { title, message } = this.props

        return (
            <Flex align="center" justify="center" className="error-page-container">
                <Flex align="center">
                    {this.renderIcon()}
                    <Flex className="error-page-text-container" vertical>
                        <h2 className="error-page-title">
                            <I18n id={title ?? MessageKey.FrontendComponentsErrorFallbackTitle} fallback="Well, that didn't work" />
                        </h2>
                        <p className="error-page-message">
                            <I18n
                                id={message ?? MessageKey.FrontendComponentsErrorFallbackMessage}
                                fallback="We ran into an error and don't know what to do with it. Please refresh the page and try again."
                            />
                        </p>
                        <Flex>{this.renderErrorButton()}</Flex>
                    </Flex>
                </Flex>
            </Flex>
        )
    }
}
