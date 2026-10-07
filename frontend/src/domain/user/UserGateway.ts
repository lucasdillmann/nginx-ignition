import ApiClient from "../../core/apiclient/ApiClient"
import UserResponse from "./model/UserResponse"
import ApiResponse from "../../core/apiclient/ApiResponse"
import UserLoginRequest from "./model/UserLoginRequest"
import UserLoginResponse from "./model/UserLoginResponse"
import UserOnboardingStatusResponse from "./model/UserOnboardingStatusResponse"
import UserRequest from "./model/UserRequest"
import PageResponse from "../../core/pagination/PageResponse"
import UserUpdatePasswordRequest from "./model/UserUpdatePasswordRequest"
import UserUpdateProfileRequest from "./model/UserUpdateProfileRequest"
import GenericCreateResponse from "../../core/common/GenericCreateResponse"
import UserTotpEnableResponse from "./model/UserTotpEnableResponse"
import TotpStatusResponse from "./model/TotpStatusResponse"
import APITokenResponse from "./model/APITokenResponse"
import APITokenRequest from "./model/APITokenRequest"
import APITokenCreatedResponse from "./model/APITokenCreatedResponse"

export default class UserGateway {
    private readonly client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/users")
    }

    async getCurrent(): Promise<ApiResponse<UserResponse>> {
        return this.client.get<UserResponse>("/current")
    }

    async getOnboardingStatus(): Promise<ApiResponse<UserOnboardingStatusResponse>> {
        return this.client.get<UserOnboardingStatusResponse>("/onboarding/status")
    }

    async finishOnboarding(request: UserRequest): Promise<ApiResponse<UserLoginResponse>> {
        return this.client.post("/onboarding/finish", request)
    }

    async login(request: UserLoginRequest): Promise<ApiResponse<UserLoginResponse>> {
        return this.client.post("/login", request)
    }

    async logout(): Promise<ApiResponse<void>> {
        return this.client.post("/logout")
    }

    async getPage(
        pageSize?: number,
        pageNumber?: number,
        searchTerms?: string,
    ): Promise<ApiResponse<PageResponse<UserResponse>>> {
        return this.client.get(undefined, undefined, { pageSize, pageNumber, searchTerms })
    }

    async getById(id: string): Promise<ApiResponse<UserResponse>> {
        return this.client.get(`/${id}`)
    }

    async putById(id: string, user: UserRequest): Promise<ApiResponse<void>> {
        return this.client.put(`/${id}`, user)
    }

    async deleteById(id: string): Promise<ApiResponse<void>> {
        return this.client.delete(`/${id}`)
    }

    async post(user: UserRequest): Promise<ApiResponse<GenericCreateResponse>> {
        return this.client.post("", user)
    }

    async updatePassword(request: UserUpdatePasswordRequest): Promise<ApiResponse<void>> {
        return this.client.post("/current/update-password", request)
    }

    async updateProfile(request: UserUpdateProfileRequest): Promise<ApiResponse<void>> {
        return this.client.put("/current", request)
    }

    async enableTotp(): Promise<ApiResponse<UserTotpEnableResponse>> {
        return this.client.post("/current/totp")
    }

    async activateTotp(code: string): Promise<ApiResponse<void>> {
        return this.client.post("/current/totp/activate", { code })
    }

    async getTotpStatus(): Promise<ApiResponse<TotpStatusResponse>> {
        return this.client.get("/current/totp")
    }

    async disableTotp(): Promise<ApiResponse<void>> {
        return this.client.delete("/current/totp")
    }

    async listTokens(
        pageSize?: number,
        pageNumber?: number,
        searchTerms?: string,
    ): Promise<ApiResponse<PageResponse<APITokenResponse>>> {
        return this.client.get("/current/tokens", undefined, { pageSize, pageNumber, searchTerms })
    }

    async createToken(request: APITokenRequest): Promise<ApiResponse<APITokenCreatedResponse>> {
        return this.client.post("/current/tokens", request)
    }

    async deleteToken(id: string): Promise<ApiResponse<void>> {
        return this.client.delete(`/current/tokens/${id}`)
    }
}
