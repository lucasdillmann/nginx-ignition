import UserGateway from "./UserGateway";
import UserResponse from "./model/UserResponse";
import UserOnboardingStatusResponse from "./model/UserOnboardingStatusResponse";
import {requireSuccessPayload} from "../../core/apiclient/ApiResponse";

export default class UserService {
    private gateway: UserGateway

    constructor() {
        this.gateway = new UserGateway()
    }

    async current(): Promise<UserResponse | undefined> {
        return this.gateway
            .getCurrent()
            .then(response => response.body)
            .catch(() => undefined)
    }

    async onboardingStatus(): Promise<UserOnboardingStatusResponse> {
        return this.gateway.getOnboardingStatus().then(requireSuccessPayload)
    }
}
