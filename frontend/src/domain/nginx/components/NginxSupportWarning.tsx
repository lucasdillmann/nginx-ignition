import NginxService from "../../nginx/NginxService"
import NginxMetadata from "../../nginx/model/NginxMetadata"
import { Alert } from "antd"
import React, { CSSProperties } from "react"

const ALERT_STYLE: CSSProperties = {
    marginBottom: 20,
}

interface NginxSupportWarningState {
    metadata?: NginxMetadata
}

export interface NginxSupportWarningMessage {
    title: string
    message: string
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
                        key={index}
                        message={title}
                        description={message}
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
