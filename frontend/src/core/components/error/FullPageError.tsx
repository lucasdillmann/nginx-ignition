import React from "react"
import { Button, Flex } from "antd"
import "./FullPageError.css"
import { ExclamationCircleFilled } from "@ant-design/icons"
import { themedModal } from "../theme/ThemedResources"

export interface FullPageErrorProps {
    title?: string
    message?: string
    error?: Error
    icon?: any
}

export default class FullPageError extends React.Component<FullPageErrorProps> {
    private openErrorDetailsModal() {
        const { error } = this.props
        themedModal().info({
            title: "Error details",
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
                Show details of what happened
            </Button>
        )
    }

    private renderIcon() {
        const { icon } = this.props
        if (icon !== undefined) return icon

        return <ExclamationCircleFilled style={{ fontSize: 48, color: "var(--nginxIgnition-colorError)" }} />
    }

    render() {
        const title = this.props.title ?? "Well, that didn't work"
        const message =
            this.props.message ??
            "We ran into an error and don't know what to do with it. Please refresh the page and try again."

        return (
            <Flex align="center" justify="center" className="error-page-container">
                <Flex align="center">
                    {this.renderIcon()}
                    <Flex className="error-page-text-container" vertical>
                        <h2 className="error-page-title">{title}</h2>
                        <p className="error-page-message">{message}</p>
                        <Flex>{this.renderErrorButton()}</Flex>
                    </Flex>
                </Flex>
            </Flex>
        )
    }
}
