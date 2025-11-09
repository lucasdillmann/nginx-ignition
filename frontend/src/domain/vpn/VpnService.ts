import VpnGateway from "./VpnGateway"
import PageResponse from "../../core/pagination/PageResponse"
import { requireNullablePayload, requireSuccessPayload, requireSuccessResponse } from "../../core/apiclient/ApiResponse"
import GenericCreateResponse from "../../core/common/GenericCreateResponse"
import VpnRequest from "./model/VpnRequest"
import VpnResponse from "./model/VpnResponse"
import AvailableDriverResponse from "./model/AvailableDriverResponse"

export default class VpnService {
    private readonly gateway: VpnGateway

    constructor() {
        this.gateway = new VpnGateway()
    }

    async getById(id: string): Promise<VpnResponse | undefined> {
        return this.gateway.getById(id).then(requireNullablePayload)
    }

    async list(
        pageSize?: number,
        pageNumber?: number,
        searchTerms?: string,
        enabledOnly?: boolean,
    ): Promise<PageResponse<VpnResponse>> {
        return this.gateway.getPage(pageSize, pageNumber, searchTerms, enabledOnly).then(requireSuccessPayload)
    }

    async delete(id: string): Promise<void> {
        return this.gateway.delete(id).then(requireSuccessResponse)
    }

    async updateById(id: string, vpn: VpnRequest): Promise<void> {
        return this.gateway.putById(id, vpn).then(requireSuccessResponse)
    }

    async create(vpn: VpnRequest): Promise<GenericCreateResponse> {
        return this.gateway.post(vpn).then(requireSuccessPayload)
    }

    async availableDrivers(): Promise<AvailableDriverResponse[]> {
        return this.gateway.getAvailableDrivers().then(requireSuccessPayload)
    }
}
