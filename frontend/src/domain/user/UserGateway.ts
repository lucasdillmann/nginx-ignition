import ApiClient from "../../core/apiclient/ApiClient"
import UserResponse from "./model/UserResponse"
import ApiResponse from "../../core/apiclient/ApiResponse"
import UserLoginRequest from "./model/UserLoginRequest"
import UserLoginResponse from "./model/UserLoginResponse"
import UserOnboardingStatusResponse from "./model/UserOnboardingStatusResponse"
import UserRequest from "./model/UserRequest"
import PageResponse from "../../core/pagination/PageResponse"
import UserUpdatePasswordRequest from "./model/UserUpdatePasswordRequest"

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

    async post(user: UserRequest): Promise<ApiResponse<void>> {
        return this.client.post("", user)
    }

    async updatePassword(request: UserUpdatePasswordRequest): Promise<ApiResponse<void>> {
        return this.client.post("/current/update-password", request)
    }
}
