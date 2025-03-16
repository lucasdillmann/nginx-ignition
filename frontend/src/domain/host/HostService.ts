import HostGateway from "./HostGateway"
import PageResponse from "../../core/pagination/PageResponse"
import HostResponse from "./model/HostResponse"
import { requireNullablePayload, requireSuccessPayload, requireSuccessResponse } from "../../core/apiclient/ApiResponse"
import HostRequest from "./model/HostRequest"
import GenericCreateResponse from "../../core/common/GenericCreateResponse"

export default class HostService {
    private readonly gateway: HostGateway

    constructor() {
        this.gateway = new HostGateway()
    }

    async getById(id: string): Promise<HostResponse | undefined> {
        return this.gateway.getById(id).then(requireNullablePayload)
    }

    async list(pageSize?: number, pageNumber?: number, searchTerms?: string): Promise<PageResponse<HostResponse>> {
        return this.gateway.getPage(pageSize, pageNumber, searchTerms).then(requireSuccessPayload)
    }

    async delete(id: string): Promise<void> {
        return this.gateway.delete(id).then(requireSuccessResponse)
    }

    async toggleEnabled(id: string): Promise<void> {
        return this.gateway.toggleEnabled(id).then(requireSuccessResponse)
    }

    async logs(id: string, type: string, lines: number): Promise<string[]> {
        return this.gateway.getLogs(id, type, lines).then(requireSuccessPayload)
    }

    async updateById(id: string, host: HostRequest): Promise<void> {
        return this.gateway.putById(id, host).then(requireSuccessResponse)
    }

    async create(host: HostRequest): Promise<GenericCreateResponse> {
        return this.gateway.post(host).then(requireSuccessPayload)
    }
}
