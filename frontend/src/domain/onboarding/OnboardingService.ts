import UserGateway from "../user/UserGateway"
import { requireSuccessPayload } from "../../core/apiclient/ApiResponse"
import AuthenticationService from "../../core/authentication/AuthenticationService"
import UserRequest from "../user/model/UserRequest"
import { UserAccessLevel } from "../user/model/UserAccessLevel"

export default class OnboardingService {
    private readonly gateway: UserGateway

    constructor() {
        this.gateway = new UserGateway()
    }

    async finish(name: string, username: string, password: string): Promise<any> {
        const request: UserRequest = {
            name: name ?? "",
            username: username ?? "",
            password: password ?? "",
            enabled: true,
            permissions: {
                hosts: UserAccessLevel.READ_WRITE,
                streams: UserAccessLevel.READ_WRITE,
                certificates: UserAccessLevel.READ_WRITE,
                logs: UserAccessLevel.READ_ONLY,
                integrations: UserAccessLevel.READ_WRITE,
                accessLists: UserAccessLevel.READ_WRITE,
                settings: UserAccessLevel.READ_WRITE,
                users: UserAccessLevel.READ_WRITE,
                nginxServer: UserAccessLevel.READ_WRITE,
            },
        }

        return this.gateway
            .finishOnboarding(request)
            .then(requireSuccessPayload)
            .then(response => {
                AuthenticationService.setToken(response.token)
            })
    }
}
