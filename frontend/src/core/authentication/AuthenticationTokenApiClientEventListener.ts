import ApiClientEventListener from "../apiclient/event/ApiClientEventListener"
import ApiResponse from "../apiclient/ApiResponse"
import AuthenticationService from "./AuthenticationService"

export default class AuthenticationTokenApiClientEventListener implements ApiClientEventListener {
    handleRequest(request: RequestInit): void {
        const token = AuthenticationService.getCurrentToken()
        if (token == null) return

        const headers = request.headers as { [key: string]: string }
        headers["Authorization"] = `Bearer ${token}`
    }

    handleResponse(_: RequestInit, response: ApiResponse<any>): void {
        const updatedAuthorization = response.headers.find(header => header.key.toLowerCase() === "authorization")

        if (updatedAuthorization != null) {
            const newToken = updatedAuthorization.value.replace("Bearer ", "")
            AuthenticationService.setToken(newToken)
        }
    }
}
