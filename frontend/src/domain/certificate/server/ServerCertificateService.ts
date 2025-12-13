import ServerCertificateGateway from "./ServerCertificateGateway"
import PageResponse from "../../../core/pagination/PageResponse"
import {
    requireNullablePayload,
    requireSuccessPayload,
    requireSuccessResponse,
    UnexpectedResponseError,
} from "../../../core/apiclient/ApiResponse"
import { ServerCertificateResponse } from "./model/ServerCertificateResponse"
import { RenewServerCertificateResponse } from "./model/RenewServerCertificateResponse"
import AvailableProviderResponse from "./model/AvailableProviderResponse"
import { IssueServerCertificateRequest } from "./model/IssueServerCertificateRequest"
import { IssueServerCertificateResponse } from "./model/IssueServerCertificateResponse"

export default class ServerCertificateService {
    private readonly gateway: ServerCertificateGateway

    constructor() {
        this.gateway = new ServerCertificateGateway()
    }

    async list(
        pageSize?: number,
        pageNumber?: number,
        searchTerms?: string,
    ): Promise<PageResponse<ServerCertificateResponse>> {
        return this.gateway.getPage(pageSize, pageNumber, searchTerms).then(requireSuccessPayload)
    }

    async delete(id: string): Promise<void> {
        return this.gateway.delete(id).then(requireSuccessResponse)
    }

    async renew(id: string): Promise<RenewServerCertificateResponse> {
        return this.gateway.renew(id).then(requireSuccessPayload)
    }

    async availableProviders(): Promise<AvailableProviderResponse[]> {
        return this.gateway.getAvailableProviders().then(requireSuccessPayload)
    }

    async issue(certificate: IssueServerCertificateRequest): Promise<IssueServerCertificateResponse> {
        return this.gateway.issue(certificate).then(response => {
            if (response.body?.success !== undefined) return response.body
            else throw new UnexpectedResponseError(response)
        })
    }

    async getById(id: string): Promise<ServerCertificateResponse | undefined> {
        return this.gateway.getById(id).then(requireNullablePayload)
    }
}
