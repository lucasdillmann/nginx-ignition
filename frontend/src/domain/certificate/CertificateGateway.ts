import ApiClient from "../../core/apiclient/ApiClient"
import ApiResponse from "../../core/apiclient/ApiResponse"
import PageResponse from "../../core/pagination/PageResponse"
import { CertificateResponse } from "./model/CertificateResponse"
import AvailableProviderResponse from "./model/AvailableProviderResponse"
import { IssueCertificateRequest } from "./model/IssueCertificateRequest"
import { SseEvent } from "../../core/apiclient/SseReader"

export default class CertificateGateway {
    private readonly client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/certificates")
    }

    async getPage(
        pageSize?: number,
        pageNumber?: number,
        searchTerms?: string,
    ): Promise<ApiResponse<PageResponse<CertificateResponse>>> {
        return this.client.get(undefined, undefined, { pageSize, pageNumber, searchTerms })
    }

    async delete(id: string): Promise<ApiResponse<void>> {
        return this.client.delete(`/${id}`)
    }

    async issueStream(
        certificate: IssueCertificateRequest,
        onEvent: (event: SseEvent) => void,
        signal?: AbortSignal,
    ): Promise<void> {
        await this.client.postStream("/issue", certificate, onEvent, signal)
    }

    async renewStream(id: string, onEvent: (event: SseEvent) => void, signal?: AbortSignal): Promise<void> {
        await this.client.postStream(`/${id}/renew`, undefined, onEvent, signal)
    }

    async getAvailableProviders(): Promise<ApiResponse<AvailableProviderResponse[]>> {
        return this.client.get("/available-providers")
    }

    async getById(id: string): Promise<ApiResponse<CertificateResponse>> {
        return this.client.get(`/${id}`)
    }
}
