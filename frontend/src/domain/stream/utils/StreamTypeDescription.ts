import { StreamType } from "../model/StreamRequest"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"

const StreamTypeDescription = {
    [StreamType.SIMPLE]: MessageKey.FrontendStreamUtilsSimple,
    [StreamType.SNI_ROUTER]: MessageKey.FrontendStreamUtilsDomainRouter,
}

export default StreamTypeDescription
