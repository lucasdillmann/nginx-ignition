import UserGateway from "./UserGateway";
import UserResponse from "./model/UserResponse";
import UserOnboardingStatusResponse from "./model/UserOnboardingStatusResponse";
import {requireSuccessPayload, requireSuccessResponse} from "../../core/apiclient/ApiResponse";
import PageResponse from "../../core/pagination/PageResponse";

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

    async list(pageSize?: number, pageNumber?: number): Promise<PageResponse<UserResponse>> {
        return this.gateway.getPage(pageSize, pageNumber).then(requireSuccessPayload)
    }

    async delete(id: string): Promise<void> {
        return this.gateway.deleteById(id).then(requireSuccessResponse)
    }
}
