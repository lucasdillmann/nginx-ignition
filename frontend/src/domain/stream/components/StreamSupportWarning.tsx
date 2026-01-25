import NginxMetadata from "../../nginx/model/NginxMetadata"
import NginxSupportWarning, { NginxSupportWarningMessage } from "../../nginx/components/NginxSupportWarning"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"

export default class StreamSupportWarning extends NginxSupportWarning {
    constructor(props: any) {
        super(props)
    }

    getWarningMessages(metadata: NginxMetadata): NginxSupportWarningMessage[] {
        const output: NginxSupportWarningMessage[] = []

        if (!metadata.availableSupport.streams)
            output.push({
                title: MessageKey.FrontendStreamComponentsSupportwarningStreamsTitle,
                message: MessageKey.FrontendStreamComponentsSupportwarningStreamsMessage,
            })

        if (!metadata.availableSupport.tlsSni)
            output.push({
                title: MessageKey.FrontendStreamComponentsSupportwarningSniTitle,
                message: MessageKey.FrontendStreamComponentsSupportwarningSniMessage,
            })

        return output
    }
}
