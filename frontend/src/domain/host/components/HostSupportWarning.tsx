import NginxService from "../../nginx/NginxService"
import NginxMetadata, { NginxSupportType } from "../../nginx/model/NginxMetadata"
import { Alert } from "antd"
import If from "../../../core/components/flowcontrol/If"
import React, { CSSProperties } from "react"

const ALERT_STYLE: CSSProperties = {
    marginBottom: 20,
}

interface HostSupportWarningState {
    metadata?: NginxMetadata
}

export default class HostSupportWarning extends React.Component<any, HostSupportWarningState> {
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
            <If condition={availableSupport.runCode == NginxSupportType.NONE}>
                <Alert
                    message="Support for code execution is not available"
                    description={
                        "The nginx server being used by nginx ignition does not support the Lua and/or " +
                        "JavaScript modules, which are both required to enable the use of code execution in the hosts' " +
                        "routes. You can still manage the hosts, but nginx will fail to start if any code execution " +
                        "route is enabled. Please contact your nginx administrator to enable the Lua/JS modules."
                    }
                    type="warning"
                    style={ALERT_STYLE}
                    showIcon
                    closable
                />
            </If>
        )
    }
}
