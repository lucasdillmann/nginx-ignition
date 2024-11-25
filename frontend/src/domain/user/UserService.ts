import UserGateway from "./UserGateway";
import UserResponse from "./model/UserResponse";
import UserOnboardingStatusResponse from "./model/UserOnboardingStatusResponse";
import {requireNullablePayload, requireSuccessPayload, requireSuccessResponse} from "../../core/apiclient/ApiResponse";
import PageResponse from "../../core/pagination/PageResponse";
import UserRequest from "./model/UserRequest";

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

    async getById(id: string): Promise<UserResponse | undefined> {
        return this.gateway.getById(id).then(requireNullablePayload)
    }

    async updateById(id: string, user: UserRequest): Promise<void> {
        return this.gateway.putById(id, user).then(requireSuccessResponse)
    }

    async create(user: UserRequest): Promise<void> {
        return this.gateway.post(user).then(requireSuccessResponse)
    }
}
