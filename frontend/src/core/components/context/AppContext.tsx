import React from "react"
import UserService from "../../../domain/user/UserService"
import UserResponse from "../../../domain/user/model/UserResponse"
import UserOnboardingStatusResponse from "../../../domain/user/model/UserOnboardingStatusResponse"

export interface AppContextData {
    user?: UserResponse
    onboardingStatus: UserOnboardingStatusResponse
}

export async function loadAppContextData(): Promise<AppContextData> {
    const service = new UserService()
    const user = service.current()
    const onboardingStatus = service.onboardingStatus()

    return Promise.all([user, onboardingStatus]).then(([user, onboardingStatus]) => ({
        user,
        onboardingStatus,
    }))
}

const AppContext = React.createContext<AppContextData>({
    user: undefined,
    onboardingStatus: {
        finished: true,
    },
})
export default AppContext
