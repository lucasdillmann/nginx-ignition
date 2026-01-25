import React from "react"
import { I18n } from "../../i18n/I18n"
import MessageKey from "../../i18n/model/MessageKey.generated"

export default {
    yesNo: (value: boolean) => <I18n id={value ? MessageKey.CommonYes : MessageKey.CommonNo} />,
}
