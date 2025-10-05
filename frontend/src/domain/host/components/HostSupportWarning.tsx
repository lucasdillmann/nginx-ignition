import NginxMetadata, { NginxSupportType } from "../../nginx/model/NginxMetadata"
import NginxSupportWarning, { NginxSupportWarningMessage } from "../../nginx/components/NginxSupportWarning"

export default class HostSupportWarning extends NginxSupportWarning {
    constructor(props: any) {
        super(props)
    }

    getWarningMessages(metadata: NginxMetadata): NginxSupportWarningMessage[] {
        if (metadata.availableSupport.runCode != NginxSupportType.NONE) return []

        return [
            {
                title: "Support for code execution is not available",
                message:
                    "The nginx server being used by nginx ignition does not support the Lua and/or " +
                    "JavaScript modules, which are both required to enable the use of code execution in the hosts' " +
                    "routes. You can still manage the hosts, but nginx will fail to start if any code execution " +
                    "route is enabled. Please contact your nginx administrator to enable the Lua/JS modules.",
            },
        ]
    }
}
