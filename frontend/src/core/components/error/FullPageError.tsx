import React from "react";
import {Button, Flex, Modal} from "antd";
import "./FullPageError.css"
import {ExclamationCircleFilled} from "@ant-design/icons";

export interface FullPageErrorProps {
    title?: string
    message?: string
    error?: Error
}

export default class FullPageError extends React.Component<FullPageErrorProps> {
    private openErrorDetailsModal() {
        const {error} = this.props
        Modal.info({
            title: "Error details",
            type: "info",
            width: 1000,
            content: (
                <pre>
                    {error?.stack || error?.message || typeof error}
                </pre>
            ),
        })
    }

    private renderErrorButton() {
        const {error} = this.props
        if (error === undefined)
            return null

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

    render() {
        const title = this.props.title ?? "Well, that didn't work"
        const message =
            this.props.message ?? "We've found an error and don't know what to do with it. " +
            "Please refresh the page and try again."

        return (
            <Flex align="center" justify="center" className="error-page-container">
                <Flex align="center">
                    <ExclamationCircleFilled style={{fontSize: 48, color: "red"}} />
                    <Flex className="error-page-text-container" vertical>
                        <h2 className="error-page-title">{title}</h2>
                        <p className="error-page-message">{message}</p>
                        <Flex>
                            {this.renderErrorButton()}
                        </Flex>
                    </Flex>
                </Flex>
            </Flex>
        )
    }
}
