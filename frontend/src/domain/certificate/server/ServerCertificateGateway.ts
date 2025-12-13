import ApiClient from "../../../core/apiclient/ApiClient"
import ApiResponse from "../../../core/apiclient/ApiResponse"
import PageResponse from "../../../core/pagination/PageResponse"
import { RenewServerCertificateResponse } from "./model/RenewServerCertificateResponse"
import { ServerCertificateResponse } from "./model/ServerCertificateResponse"
import AvailableProviderResponse from "./model/AvailableProviderResponse"
import { IssueServerCertificateRequest } from "./model/IssueServerCertificateRequest"
import { IssueServerCertificateResponse } from "./model/IssueServerCertificateResponse"

export default class ServerCertificateGateway {
    private readonly client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/certificates/server")
    }

    async getPage(
        pageSize?: number,
        pageNumber?: number,
        searchTerms?: string,
    ): Promise<ApiResponse<PageResponse<ServerCertificateResponse>>> {
        return this.client.get(undefined, undefined, { pageSize, pageNumber, searchTerms })
    }

    async delete(id: string): Promise<ApiResponse<void>> {
        return this.client.delete(`/${id}`)
    }

    async renew(id: string): Promise<ApiResponse<RenewServerCertificateResponse>> {
        return this.client.post(`/${id}/renew`)
    }

    async getAvailableProviders(): Promise<ApiResponse<AvailableProviderResponse[]>> {
        return this.client.get("/available-providers")
    }

    async issue(certificate: IssueServerCertificateRequest): Promise<ApiResponse<IssueServerCertificateResponse>> {
        return this.client.post("/issue", certificate)
    }

    async getById(id: string): Promise<ApiResponse<ServerCertificateResponse>> {
        return this.client.get(`/${id}`)
    }
}
