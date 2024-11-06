import ApiResponse from "./ApiResponse"
import Header from "./Header"
import ApiClientEventDispatcher from "./event/ApiClientEventDispatcher";

export default class ApiClient {
    private readonly basePath: string

    constructor(basePath?: string) {
        this.basePath = basePath || ""
    }

    async get<T>(
        path?: string,
        headers?: Header[],
        queryParams?: Map<string, string>,
    ): Promise<ApiResponse<T>> {
        const request = await this.buildRequest("GET", headers)
        const fullPath = this.buildFullPath(path, queryParams)
        return await this.executeRequest(fullPath, request)
    }

    async delete<T>(
        path?: string,
        headers?: Header[],
        queryParams?: Map<string, string>,
    ): Promise<ApiResponse<T>> {
        const request = await this.buildRequest("DELETE", headers)
        const fullPath = this.buildFullPath(path, queryParams)
        return await this.executeRequest(fullPath, request)
    }

    async post<I, O>(
        path?: string,
        payload?: I,
        headers?: Header[],
        queryParams?: Map<string, string>
    ): Promise<ApiResponse<O>> {
        const request = await this.buildRequest("POST", headers, payload)
        const fullPath = this.buildFullPath(path, queryParams)
        return await this.executeRequest(fullPath, request)
    }

    async put<I, O>(
        path?: string,
        payload?: I,
        headers?: Header[],
        queryParams?: Map<string, string>
    ): Promise<ApiResponse<O>> {
        const request = await this.buildRequest("PUT", headers, payload)
        const fullPath = this.buildFullPath(path, queryParams)
        return await this.executeRequest(fullPath, request)
    }

    private async buildRequest<T>(method: string, headers?: Header[], payload?: T): Promise<RequestInit> {
        const requestHeaders = new Headers()
        requestHeaders.append("Accept", "application/json")
        requestHeaders.append("Content-type", "application/json")

        headers?.forEach(({ key, value }) => {
            requestHeaders.append(key, value)
        })

        const body = JSON.stringify(payload)

        return {
            method,
            headers: requestHeaders,
            body,
        }
    }

    private buildFullPath(path?: string, queryParams?: Map<string, string>): string {
        const queryString = !queryParams
            ? ""
            : "?" +
              Array.from(queryParams.entries())
                  .map(([key, value]) => `${key}=${value}`)
                  .join("&")
        if (path == null) return this.basePath
        else return `${this.basePath}${path}${queryString}`
    }

    private async executeRequest<T>(path: string, request: RequestInit): Promise<ApiResponse<T>> {
        ApiClientEventDispatcher.notifyRequest(request)

        const response = await fetch(path, request)
        const headers: Header[] = Array
            .from(response.headers.entries())
            .map(([key, value]) => ({ key, value }))

        let body
        try {
            body = await response.json()
        } catch (e) {
        }

        const apiResponse: ApiResponse<T> = {
            statusCode: response.status,
            body,
            headers,
        }

        ApiClientEventDispatcher.notifyResponse(request, apiResponse)
        return apiResponse;
    }
}
