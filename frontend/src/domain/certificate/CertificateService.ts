import CertificateGateway from "./CertificateGateway";
import PageResponse from "../../core/pagination/PageResponse";
import {requireSuccessPayload, requireSuccessResponse} from "../../core/apiclient/ApiResponse";
import {CertificateResponse} from "./model/CertificateResponse";
import {RenewCertificateResponse} from "./model/RenewCertificateResponse";
import AvailableProviderResponse from "./model/AvailableProviderResponse";

export default class CertificateService {
    private readonly gateway: CertificateGateway

    constructor() {
        this.gateway = new CertificateGateway()
    }

    async list(pageSize?: number, pageNumber?: number): Promise<PageResponse<CertificateResponse>> {
        return this.gateway.getPage(pageSize, pageNumber).then(requireSuccessPayload)
    }

    async delete(id: string): Promise<void> {
        return this.gateway.delete(id).then(requireSuccessResponse)
    }

    async renew(id: string): Promise<RenewCertificateResponse> {
        return this.gateway.renew(id).then(requireSuccessPayload)
    }

    async availableProviders(): Promise<AvailableProviderResponse[]> {
        return this.gateway.getAvailableProviders().then(requireSuccessPayload)
    }
}
