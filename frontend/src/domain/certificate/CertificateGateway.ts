import ApiClient from "../../core/apiclient/ApiClient"
import ApiResponse from "../../core/apiclient/ApiResponse"
import PageResponse from "../../core/pagination/PageResponse"
import { RenewCertificateResponse } from "./model/RenewCertificateResponse"
import { CertificateResponse } from "./model/CertificateResponse"
import AvailableProviderResponse from "./model/AvailableProviderResponse"
import { IssueCertificateRequest } from "./model/IssueCertificateRequest"
import { IssueCertificateResponse } from "./model/IssueCertificateResponse"

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

    async renew(id: string): Promise<ApiResponse<RenewCertificateResponse>> {
        return this.client.post(`/${id}/renew`)
    }

    async getAvailableProviders(): Promise<ApiResponse<AvailableProviderResponse[]>> {
        return this.client.get("/available-providers")
    }

    async issue(certificate: IssueCertificateRequest): Promise<ApiResponse<IssueCertificateResponse>> {
        return this.client.post("/issue", certificate)
    }

    async getById(id: string): Promise<ApiResponse<CertificateResponse>> {
        return this.client.get(`/${id}`)
    }
}
