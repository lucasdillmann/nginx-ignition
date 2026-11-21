import React from "react"
import { Button, Checkbox, Flex, Typography } from "antd"
import { CheckCircleFilled, CloseCircleFilled, ExclamationCircleFilled, LoadingOutlined } from "@ant-design/icons"
import dayjs from "dayjs"
import "./CertificateOperationModal.css"
import "../../../core/components/preloader/ModalPreloader.css"
import {
    CertificateOperationLogLine,
    CertificateOperationState,
    CertificateOperationStore,
} from "./CertificateOperationStore"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { I18n, I18nMessage, raw } from "../../../core/i18n/I18n"
import { navigateTo } from "../../../core/components/router/AppRouter"
import ReloadNginxAction from "../../nginx/actions/ReloadNginxAction"

export interface CertificateOperationModalProps {
    store: CertificateOperationStore
    onClose: () => void
}

interface CertificateOperationModalState {
    operation: CertificateOperationState
}

export default class CertificateOperationModal extends React.Component<
    CertificateOperationModalProps,
    CertificateOperationModalState
> {
    private readonly logContainerRef: React.RefObject<HTMLDivElement | null>
    private unsubscribe?: () => void

    constructor(props: CertificateOperationModalProps) {
        super(props)
        this.logContainerRef = React.createRef()
        this.state = {
            operation: props.store.getState(),
        }
    }

    componentDidMount() {
        this.unsubscribe = this.props.store.subscribe(() => {
            this.setState({ operation: this.props.store.getState() })
        })

        this.scrollToBottom()
    }

    componentDidUpdate(_: CertificateOperationModalProps, previousState: CertificateOperationModalState) {
        if (previousState.operation.logCount !== this.state.operation.logCount) this.scrollToBottom()
    }

    componentWillUnmount() {
        this.unsubscribe?.()
    }

    private scrollToBottom() {
        const container = this.logContainerRef.current
        if (container !== null) container.scrollTop = container.scrollHeight
    }

    private closeWith(action?: () => void | Promise<void>) {
        this.props.onClose()
        action?.()
    }

    private renderLogLine(line: CertificateOperationLogLine) {
        return (
            <Flex
                key={line.id}
                className={`certificate-operation-log-line certificate-operation-log-line-${line.level}`}
            >
                <span className="certificate-operation-log-time">{dayjs(line.time).format("HH:mm:ss")}</span>
                <span className={`certificate-operation-log-level certificate-operation-log-level-${line.level}`}>
                    {line.level.toUpperCase()}
                </span>
                <span className="certificate-operation-log-message">{line.message}</span>
            </Flex>
        )
    }

    private renderLogs() {
        const { logs, showDebug, truncated } = this.state.operation
        const visibleLogs = showDebug ? logs : logs.filter(line => line.level !== "debug")

        return (
            <div ref={this.logContainerRef} className="certificate-operation-log-container">
                {truncated && (
                    <div className="certificate-operation-log-notice">
                        <I18n id={MessageKey.FrontendCertificateOperationLogsTruncated} />
                    </div>
                )}
                {visibleLogs.length === 0 && (
                    <div className="certificate-operation-log-notice">
                        <I18n id={MessageKey.FrontendCertificateOperationLogsEmpty} />
                    </div>
                )}
                {visibleLogs.map(line => this.renderLogLine(line))}
            </div>
        )
    }

    private renderRunning() {
        const { kind, phase } = this.state.operation
        const message =
            phase === "canceling"
                ? MessageKey.FrontendCertificateOperationCanceling
                : kind === "issue"
                  ? MessageKey.FrontendCertificateIssuingPleaseWait
                  : MessageKey.FrontendCertificateRenewingPleaseWait

        return (
            <Flex className="certificate-operation-status" align="center">
                <LoadingOutlined className="modal-preloader-spinner" spin />
                <Typography.Text>
                    <I18n id={message} />
                </Typography.Text>
            </Flex>
        )
    }

    private renderResultBlock(icon: React.ReactNode, title: MessageKey, description?: I18nMessage) {
        return (
            <Flex className="certificate-operation-result" align="flex-start" gap="small">
                {icon}
                <Flex vertical>
                    <Typography.Text strong>
                        <I18n id={title} />
                    </Typography.Text>
                    {description !== undefined && (
                        <Typography.Text type="secondary">
                            <I18n id={description} />
                        </Typography.Text>
                    )}
                </Flex>
            </Flex>
        )
    }

    private renderResult() {
        const { kind, phase, errorReason } = this.state.operation

        if (phase === "success") {
            const title =
                kind === "issue" ? MessageKey.FrontendCertificateIssued : MessageKey.FrontendCertificateRenewed
            const description =
                kind === "issue"
                    ? MessageKey.FrontendCertificateIssueSuccessDescription
                    : MessageKey.CommonSuccessMessage

            return this.renderResultBlock(
                <CheckCircleFilled className="certificate-operation-result-icon certificate-operation-result-icon-success" />,
                title,
                description,
            )
        }

        if (phase === "canceled") {
            return this.renderResultBlock(
                <ExclamationCircleFilled className="certificate-operation-result-icon certificate-operation-result-icon-warning" />,
                MessageKey.FrontendCertificateOperationCanceled,
            )
        }

        const title =
            kind === "issue" ? MessageKey.FrontendCertificateIssueFailed : MessageKey.FrontendCertificateRenewFailed
        const description = errorReason
            ? raw(errorReason)
            : kind === "issue"
              ? MessageKey.FrontendCertificateIssueFailedDescription
              : MessageKey.CommonUnexpectedErrorTryAgain

        return this.renderResultBlock(
            <CloseCircleFilled className="certificate-operation-result-icon certificate-operation-result-icon-error" />,
            title,
            description,
        )
    }

    private renderActions() {
        const { kind, phase, certificateId } = this.state.operation

        if (phase === "running" || phase === "canceling") {
            return (
                <Button danger loading={phase === "canceling"} onClick={() => this.props.store.cancel()}>
                    <I18n id={MessageKey.CommonCancel} />
                </Button>
            )
        }

        if (phase === "success" && kind === "issue" && certificateId !== undefined) {
            return (
                <Flex gap="small">
                    <Button onClick={() => this.closeWith()}>
                        <I18n id={MessageKey.FrontendCertificateOperationClose} />
                    </Button>
                    <Button
                        type="primary"
                        onClick={() => this.closeWith(() => navigateTo(`/certificates/${certificateId}`))}
                    >
                        <I18n id={MessageKey.FrontendCertificateOperationViewCertificate} />
                    </Button>
                </Flex>
            )
        }

        if (phase === "success" && kind === "renew") {
            return (
                <Flex gap="small">
                    <Button onClick={() => this.closeWith()}>
                        <I18n id={MessageKey.FrontendCertificateOperationClose} />
                    </Button>
                    <Button type="primary" onClick={() => this.closeWith(() => ReloadNginxAction.execute())}>
                        <I18n id={MessageKey.FrontendCertificateOperationReloadNginx} />
                    </Button>
                </Flex>
            )
        }

        return (
            <Button onClick={() => this.closeWith()}>
                <I18n id={MessageKey.FrontendCertificateOperationClose} />
            </Button>
        )
    }

    private renderStatus() {
        const { phase } = this.state.operation
        if (phase === "running" || phase === "canceling") return this.renderRunning()

        return this.renderResult()
    }

    render() {
        const { showDebug } = this.state.operation

        return (
            <Flex className="certificate-operation-modal-content" vertical gap="middle">
                {this.renderStatus()}

                <Flex justify="space-between" align="center">
                    <Typography.Text strong>
                        <I18n id={MessageKey.FrontendCertificateOperationLiveLog} />
                    </Typography.Text>
                    <Checkbox
                        checked={showDebug}
                        onChange={event => this.props.store.setShowDebug(event.target.checked)}
                    >
                        <I18n id={MessageKey.FrontendCertificateOperationShowDebugLogs} />
                    </Checkbox>
                </Flex>

                {this.renderLogs()}

                <Flex justify="end">{this.renderActions()}</Flex>
            </Flex>
        )
    }
}
