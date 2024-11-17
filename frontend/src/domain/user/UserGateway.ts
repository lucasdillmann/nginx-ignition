import ApiClient from "../../core/apiclient/ApiClient";
import UserResponse from "./model/UserResponse";
import ApiResponse from "../../core/apiclient/ApiResponse";
import UserLoginRequest from "./model/UserLoginRequest";
import UserLoginResponse from "./model/UserLoginResponse";
import UserOnboardingStatusResponse from "./model/UserOnboardingStatusResponse";
import UserRequest from "./model/UserRequest";

export default class UserGateway {
    private client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/users")
    }

    getCurrent(): Promise<ApiResponse<UserResponse>> {
        return this.client.get<UserResponse>("/current")
    }

    getOnboardingStatus(): Promise<ApiResponse<UserOnboardingStatusResponse>> {
        return this.client.get<UserOnboardingStatusResponse>("/onboarding/status")
    }

    finishOnboarding(request: UserRequest): Promise<ApiResponse<UserLoginResponse>> {
        return this.client.post("/onboarding/finish", request)
    }

    login(request: UserLoginRequest): Promise<ApiResponse<UserLoginResponse>> {
        return this.client.post("/login", request)
    }
}
