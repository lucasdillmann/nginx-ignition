import React from "react"
import ApiClientEventListener from "../apiclient/event/ApiClientEventListener"
import ApiResponse from "../apiclient/ApiResponse"
import AppContext from "../components/context/AppContext"
import { navigateTo } from "../components/router/AppRouter"
import { buildLoginUrl } from "./buildLoginUrl"
import { themedModal } from "../components/theme/ThemedResources"

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
            title: "Your AFK was too long",
            content: (
                <>
                    <p>
                        Seems like you've kept the nginx ignition open but without using it for a while and your session
                        expired.
                    </p>
                    <p>
                        We need to login again to continue, but you can choose to stay here if you have some unfinished
                        action that you don't want to lose or start from scratch (you can try to finish it again after
                        logging-in in another tab).
                    </p>
                </>
            ),
            okText: "Login again",
            cancelText: "Keep me here for now",
            onOk: () => this.redirectToLogin(),
        })
    }
}
