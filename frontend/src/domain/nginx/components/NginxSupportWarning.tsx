import NginxService from "../../nginx/NginxService"
import NginxMetadata from "../../nginx/model/NginxMetadata"
import { Alert } from "antd"
import React, { CSSProperties } from "react"
import { I18n, I18nMessage } from "../../../core/i18n/I18n"

const ALERT_STYLE: CSSProperties = {
    marginBottom: 20,
}

interface NginxSupportWarningState {
    metadata?: NginxMetadata
}

export interface NginxSupportWarningMessage {
    title: I18nMessage
    message: I18nMessage
}

export default abstract class NginxSupportWarning extends React.Component<any, NginxSupportWarningState> {
    private readonly service: NginxService

    constructor(props: any) {
        super(props)
        this.service = new NginxService()
        this.state = {}
    }

    abstract getWarningMessages(metadata: NginxMetadata): NginxSupportWarningMessage[]

    componentDidMount() {
        this.service
            .getMetadata()
            .then(metadata => this.setState({ metadata }))
            .catch(() => this.setState({ metadata: undefined }))
    }

    render() {
        const { metadata } = this.state
        if (metadata == null) return null

        const messages = this.getWarningMessages(metadata)

        return (
            <>
                {messages.map(({ title, message }, index) => (
                    <Alert
                        key={`nginx-support-alert-${index}`}
                        message={<I18n id={title} />}
                        description={<I18n id={message} />}
                        type="warning"
                        style={ALERT_STYLE}
                        showIcon
                        closable
                    />
                ))}
            </>
        )
    }
}
