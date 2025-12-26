import CacheGateway from "./CacheGateway"
import { requireNullablePayload, requireSuccessPayload, requireSuccessResponse } from "../../core/apiclient/ApiResponse"
import PageResponse from "../../core/pagination/PageResponse"
import CacheRequest from "./model/CacheRequest"
import CacheResponse from "./model/CacheResponse"
import GenericCreateResponse from "../../core/common/GenericCreateResponse"

export default class CacheService {
    private readonly gateway: CacheGateway

    constructor() {
        this.gateway = new CacheGateway()
    }

    async list(
        pageSize?: number,
        pageNumber?: number,
        searchTerms?: string,
    ): Promise<PageResponse<CacheResponse>> {
        return this.gateway.getPage(pageSize, pageNumber, searchTerms).then(requireSuccessPayload)
    }

    async delete(id: string): Promise<void> {
        return this.gateway.deleteById(id).then(requireSuccessResponse)
    }

    async getById(id: string): Promise<CacheResponse | undefined> {
        return this.gateway.getById(id).then(requireNullablePayload)
    }

    async updateById(id: string, user: CacheRequest): Promise<void> {
        return this.gateway.putById(id, user).then(requireSuccessResponse)
    }

    async create(user: CacheRequest): Promise<GenericCreateResponse> {
        return this.gateway.post(user).then(requireSuccessPayload)
    }
}
