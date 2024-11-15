import ApiClient from "../../core/apiclient/ApiClient";
import UserResponse from "./model/UserResponse";
import ApiResponse from "../../core/apiclient/ApiResponse";

export default class UserGateway {
    private client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/users")
    }

    getCurrent(): Promise<ApiResponse<UserResponse>> {
        return this.client.get<UserResponse>("/current")
    }
}
