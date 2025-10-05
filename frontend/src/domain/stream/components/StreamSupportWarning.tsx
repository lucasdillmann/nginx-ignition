import NginxService from "../../nginx/NginxService"
import NginxMetadata, { NginxSupportType } from "../../nginx/model/NginxMetadata"
import { Alert } from "antd"
import If from "../../../core/components/flowcontrol/If"
import React, { CSSProperties } from "react"

const ALERT_STYLE: CSSProperties = {
    marginBottom: 20,
}

interface StreamSupportWarningState {
    metadata?: NginxMetadata
}

export default class StreamSupportWarning extends React.Component<any, StreamSupportWarningState> {
    private readonly service: NginxService

    constructor(props: any) {
        super(props)
        this.service = new NginxService()
        this.state = {}
    }

    componentDidMount() {
        this.service
            .getMetadata()
            .then(metadata => this.setState({ metadata }))
            .catch(() => this.setState({ metadata: undefined }))
    }

    render() {
        const { metadata } = this.state
        if (metadata == null) return null

        const { availableSupport } = metadata

        return (
            <>
                <If condition={availableSupport.streams == NginxSupportType.NONE}>
                    <Alert
                        message="Support for streams is not available"
                        description={
                            "The nginx server being used by nginx ignition does not support the streams module. You " +
                            "can still manage the streams but nginx ignition will not be able to start the nginx " +
                            "server. Please contact your nginx administrator to enable the streams module."
                        }
                        type="warning"
                        style={ALERT_STYLE}
                        showIcon
                        closable
                    />
                </If>
                <If condition={!availableSupport.tlsSni}>
                    <Alert
                        message="Support for TLS SNI is not available"
                        description={
                            "The nginx server being used by nginx ignition does not support TLS SNI (Server Name " +
                            "Indication), a required feature for the domain-based routed streams to work. You can " +
                            "still manage the streams, but nginx will fail to start if any domain-based stream is " +
                            "enabled. Please contact your nginx administrator to enable the SNI support."
                        }
                        type="warning"
                        style={ALERT_STYLE}
                        showIcon
                        closable
                    />
                </If>
            </>
        )
    }
}
