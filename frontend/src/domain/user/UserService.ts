import UserGateway from "./UserGateway"
import UserResponse from "./model/UserResponse"
import UserOnboardingStatusResponse from "./model/UserOnboardingStatusResponse"
import { requireNullablePayload, requireSuccessPayload, requireSuccessResponse } from "../../core/apiclient/ApiResponse"
import PageResponse from "../../core/pagination/PageResponse"
import UserRequest from "./model/UserRequest"
import LoginOutcome from "./model/LoginOutcome"
import UserLoginRequest from "./model/UserLoginRequest"
import AuthenticationService from "../../core/authentication/AuthenticationService"
import UserUpdatePasswordRequest from "./model/UserUpdatePasswordRequest"
import GenericCreateResponse from "../../core/common/GenericCreateResponse"

export default class UserService {
    private readonly gateway: UserGateway

    constructor() {
        this.gateway = new UserGateway()
    }

    async login(username: string, password: string, totp?: string): Promise<LoginOutcome> {
        const request: UserLoginRequest = { username, password, totp }
        const response = await this.gateway.login(request)

        if (response.statusCode >= 200 && response.statusCode <= 299 && response.body?.token) {
            AuthenticationService.setToken(response.body.token)
            return LoginOutcome.SUCCESS
        }

        return (response.body as any)?.reason ?? LoginOutcome.FAILURE
    }

    async logout(): Promise<void> {
        return this.gateway
            .logout()
            .catch(() => undefined)
            .then(() => AuthenticationService.deleteToken())
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

    async list(pageSize?: number, pageNumber?: number, searchTerms?: string): Promise<PageResponse<UserResponse>> {
        return this.gateway.getPage(pageSize, pageNumber, searchTerms).then(requireSuccessPayload)
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

    async create(user: UserRequest): Promise<GenericCreateResponse> {
        return this.gateway.post(user).then(requireSuccessPayload)
    }

    async changePassword(request: UserUpdatePasswordRequest): Promise<void> {
        return this.gateway.updatePassword(request).then(requireSuccessResponse)
    }
}
