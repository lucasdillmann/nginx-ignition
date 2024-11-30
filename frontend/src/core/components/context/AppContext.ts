import UserService from "../../../domain/user/UserService"
import UserResponse from "../../../domain/user/model/UserResponse"
import UserOnboardingStatusResponse from "../../../domain/user/model/UserOnboardingStatusResponse"
import ContextHolder from "../../context/ContextHolder"

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

// eslint-disable-next-line import/no-anonymous-default-export
export default new ContextHolder<AppContextData>({
    user: undefined,
    onboardingStatus: {
        finished: true,
    },
})
