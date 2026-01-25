import React from "react"
import { themedModal } from "../theme/ThemedResources"
import { I18n } from "../../i18n/I18n"
import MessageKey from "../../i18n/model/MessageKey.generated"

class AccessDeniedModal {
    show() {
        themedModal().error({
            title: <I18n id={MessageKey.FrontendComponentsAccesscontrolAccessDeniedTitle} />,
            content: <I18n id={MessageKey.FrontendComponentsAccesscontrolActionDeniedDescription} />,
        })
    }
}

export default new AccessDeniedModal()
