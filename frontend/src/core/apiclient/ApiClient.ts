import ApiResponse, { UnexpectedResponseError } from "./ApiResponse"
import Header from "./Header"
import ApiClientEventDispatcher from "./event/ApiClientEventDispatcher"
import I18nContext from "../i18n/I18nContext"
import { readEventStream, SseEvent } from "./SseReader"

export default class ApiClient {
    private readonly basePath: string

    constructor(basePath?: string) {
        this.basePath = basePath ?? ""
    }

    async get<T>(
        path?: string,
        headers?: Header[],
        queryParams?: { [key: string]: any },
        raw: boolean = false,
    ): Promise<ApiResponse<T>> {
        const request = await this.buildRequest("GET", headers)
        const fullPath = this.buildFullPath(path, queryParams)
        return await this.executeRequest(fullPath, request, raw)
    }

    async delete<T>(path?: string, headers?: Header[], queryParams?: { [key: string]: any }): Promise<ApiResponse<T>> {
        const request = await this.buildRequest("DELETE", headers)
        const fullPath = this.buildFullPath(path, queryParams)
        return await this.executeRequest(fullPath, request)
    }

    async post<I, O>(
        path?: string,
        payload?: I,
        headers?: Header[],
        queryParams?: { [key: string]: any },
    ): Promise<ApiResponse<O>> {
        const request = await this.buildRequest("POST", headers, payload)
        const fullPath = this.buildFullPath(path, queryParams)
        return await this.executeRequest(fullPath, request)
    }

    async put<I, O>(
        path?: string,
        payload?: I,
        headers?: Header[],
        queryParams?: { [key: string]: any },
    ): Promise<ApiResponse<O>> {
        const request = await this.buildRequest("PUT", headers, payload)
        const fullPath = this.buildFullPath(path, queryParams)
        return await this.executeRequest(fullPath, request)
    }

    async postStream<I>(
        path: string,
        payload: I | undefined,
        onEvent: (event: SseEvent) => void,
        signal?: AbortSignal,
        headers?: Header[],
    ): Promise<ApiResponse<void>> {
        const request = await this.buildRequest("POST", headers, payload)
        const requestHeaders = request.headers as { [key: string]: string }
        requestHeaders["Accept"] = "text/event-stream"
        request.signal = signal

        ApiClientEventDispatcher.notifyRequest(request)

        const fullPath = this.buildFullPath(path)
        const response = await fetch(fullPath, request)
        const responseHeaders: Header[] = Array.from(response.headers.entries()).map(([key, value]) => ({ key, value }))
        const contentType = response.headers.get("content-type") ?? ""

        if (!response.ok || !contentType.includes("text/event-stream")) {
            const body = await this.parseResponseBody(response)
            const apiResponse: ApiResponse<any> = {
                raw: response,
                statusCode: response.status,
                body,
                headers: responseHeaders,
            }

            ApiClientEventDispatcher.notifyResponse(request, apiResponse)
            throw new UnexpectedResponseError(apiResponse)
        }

        const apiResponse: ApiResponse<void> = {
            raw: response,
            statusCode: response.status,
            headers: responseHeaders,
        }

        ApiClientEventDispatcher.notifyResponse(request, apiResponse)
        await readEventStream(response, onEvent)

        return apiResponse
    }

    private async buildRequest<T>(method: string, headers?: Header[], payload?: T): Promise<RequestInit> {
        const { currentLanguage, defaultLanguage } = I18nContext.get()
        const requestHeaders: { [key: string]: string } = {
            Accept: "application/json",
            "Accept-Language": currentLanguage ?? defaultLanguage,
            "Content-type": "application/json",
        }

        headers?.forEach(({ key, value }) => {
            requestHeaders[key] = value
        })

        const body = JSON.stringify(payload)

        return {
            method,
            headers: requestHeaders,
            body,
        }
    }

    private buildFullPath(path?: string, queryParams?: { [key: string]: any }): string {
        const queryStringValues = Array.from(Object.entries(queryParams ?? {}))
            .filter(([_, value]) => value !== undefined)
            .map(([key, value]) => `${key}=${value}`)
            .join("&")
        const queryString = queryStringValues.length > 0 ? `?${queryStringValues}` : ""

        return `${this.basePath}${path ?? ""}${queryString}`
    }

    private async executeRequest<T>(path: string, request: RequestInit, raw: boolean = false): Promise<ApiResponse<T>> {
        ApiClientEventDispatcher.notifyRequest(request)

        const response = await fetch(path, request)
        const headers: Header[] = Array.from(response.headers.entries()).map(([key, value]) => ({ key, value }))

        let body
        if (!raw) body = await this.parseResponseBody(response)

        const apiResponse: ApiResponse<T> = {
            raw: response,
            statusCode: response.status,
            body,
            headers,
        }

        ApiClientEventDispatcher.notifyResponse(request, apiResponse)
        return apiResponse
    }

    private async parseResponseBody(response: Response): Promise<any> {
        try {
            return await response.json()
        } catch (_) {
            return undefined
        }
    }
}
