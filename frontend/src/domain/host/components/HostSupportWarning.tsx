import NginxMetadata, { NginxSupportType } from "../../nginx/model/NginxMetadata"
import NginxSupportWarning, { NginxSupportWarningMessage } from "../../nginx/components/NginxSupportWarning"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"

export default class HostSupportWarning extends NginxSupportWarning {
    constructor(props: any) {
        super(props)
    }

    getWarningMessages(metadata: NginxMetadata): NginxSupportWarningMessage[] {
        if (metadata.availableSupport.runCode != NginxSupportType.NONE) return []

        return [
            {
                title: MessageKey.FrontendHostComponentsHostsupportwarningTitle,
                message: MessageKey.FrontendHostComponentsHostsupportwarningDescription,
            },
        ]
    }
}
