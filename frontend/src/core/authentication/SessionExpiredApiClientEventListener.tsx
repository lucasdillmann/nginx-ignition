import React from "react"
import ApiClientEventListener from "../apiclient/event/ApiClientEventListener"
import ApiResponse from "../apiclient/ApiResponse"
import AppContext from "../components/context/AppContext"
import { navigateTo } from "../components/router/AppRouter"
import { buildLoginUrl } from "./buildLoginUrl"
import { themedModal } from "../components/theme/ThemedResources"
import { I18n } from "../i18n/I18n"
import MessageKey from "../i18n/model/MessageKey.generated"

export default class SessionExpiredApiClientEventListener implements ApiClientEventListener {
    private alreadyShown: boolean

    constructor() {
        this.alreadyShown = false
    }

    private redirectToLogin() {
        AppContext.get().user = undefined
        navigateTo(buildLoginUrl())
    }

    handleRequest(_: RequestInit): void {
        // NO-OP
    }

    handleResponse(_: RequestInit, response: ApiResponse<any>): void {
        const currentUser = AppContext.get().user
        if (currentUser === undefined || response.statusCode !== 401 || this.alreadyShown) return

        this.alreadyShown = true

        themedModal().confirm({
            title: <I18n id={MessageKey.FrontendAuthenticationSessionExpiredTitle} />,
            content: (
                <>
                    <p>
                        <I18n id={MessageKey.FrontendAuthenticationSessionExpiredDescription1} />
                    </p>
                    <p>
                        <I18n id={MessageKey.FrontendAuthenticationSessionExpiredDescription2} />
                    </p>
                </>
            ),
            okText: <I18n id={MessageKey.FrontendAuthenticationSessionExpiredLogin} />,
            cancelText: <I18n id={MessageKey.FrontendAuthenticationSessionExpiredStay} />,
            onOk: () => this.redirectToLogin(),
        })
    }
}
