import NginxMetadata from "../../nginx/model/NginxMetadata"
import NginxSupportWarning, { NginxSupportWarningMessage } from "../../nginx/components/NginxSupportWarning"

export default class StreamSupportWarning extends NginxSupportWarning {
    constructor(props: any) {
        super(props)
    }

    getWarningMessages(metadata: NginxMetadata): NginxSupportWarningMessage[] {
        const output: NginxSupportWarningMessage[] = []

        if (!metadata.availableSupport.streams)
            output.push({
                title: "Support for streams is not available",
                message:
                    "The nginx server being used by nginx ignition does not support the streams module. You " +
                    "can still manage the streams but nginx ignition will not be able to start the nginx " +
                    "server. Please contact your nginx administrator to enable the streams module.",
            })

        if (!metadata.availableSupport.tlsSni)
            output.push({
                title: "Support for TLS SNI is not available",
                message:
                    "The nginx server being used by nginx ignition does not support TLS SNI (Server Name " +
                    "Indication), a required feature for the domain-based routed streams to work. You can " +
                    "still manage the streams, but nginx will fail to start if any domain-based stream is " +
                    "enabled. Please contact your nginx administrator to enable the SNI support.",
            })

        return output
    }
}
