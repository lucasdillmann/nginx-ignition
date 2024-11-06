import ApiResponse from "../ApiResponse";

export default interface ApiClientEventListener {
    handleRequest(request: RequestInit): void
    handleResponse(request: RequestInit, response: ApiResponse<any>): void
}
