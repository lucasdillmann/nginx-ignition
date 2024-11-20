import UserGateway from "../user/UserGateway";
import {requireSuccessPayload} from "../../core/apiclient/ApiResponse";
import AuthenticationService from "../../core/authentication/AuthenticationService";
import UserRequest from "../user/model/UserRequest";
import {UserRole} from "../user/model/UserRole";

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
            role: UserRole.ADMIN,
        }

        return this.gateway
            .finishOnboarding(request)
            .then(requireSuccessPayload)
            .then((response) => {
                AuthenticationService.setToken(response.token)
            })
    }
}
