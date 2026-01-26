import UserService from "../../../domain/user/UserService"
import UserResponse from "../../../domain/user/model/UserResponse"
import UserOnboardingStatusResponse from "../../../domain/user/model/UserOnboardingStatusResponse"
import ContextHolder from "../../context/ContextHolder"
import Configuration from "../../../domain/configuration/Configuration"
import ConfigurationService from "../../../domain/configuration/ConfigurationService"
import AppContainer from "../../../domain/AppContainer"

export interface AppContextData {
    container?: AppContainer
    user?: UserResponse
    onboardingStatus: UserOnboardingStatusResponse
    configuration: Configuration
}

export async function loadAppContextData(): Promise<AppContextData> {
    const userService = new UserService()
    const configurationService = new ConfigurationService()
    const user = userService.current()
    const onboardingStatus = userService.onboardingStatus()
    const configuration = configurationService.get()

    return Promise.all([user, onboardingStatus, configuration]).then(([user, onboardingStatus, configuration]) => ({
        user,
        onboardingStatus,
        configuration,
    }))
}

export default new ContextHolder<AppContextData>({
    user: undefined,
    onboardingStatus: {
        finished: true,
    },
    configuration: {
        version: {},
        codeEditor: {
            apiKey: undefined,
        },
    },
})
