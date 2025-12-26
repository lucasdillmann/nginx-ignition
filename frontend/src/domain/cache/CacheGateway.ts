import ApiClient from "../../core/apiclient/ApiClient"
import ApiResponse from "../../core/apiclient/ApiResponse"
import PageResponse from "../../core/pagination/PageResponse"
import CacheResponse from "./model/CacheResponse"
import CacheRequest from "./model/CacheRequest"
import GenericCreateResponse from "../../core/common/GenericCreateResponse"

export default class CacheGateway {
    private readonly client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/caches")
    }

    async getPage(
        pageSize?: number,
        pageNumber?: number,
        searchTerms?: string,
    ): Promise<ApiResponse<PageResponse<CacheResponse>>> {
        return this.client.get(undefined, undefined, { pageSize, pageNumber, searchTerms })
    }

    async getById(id: string): Promise<ApiResponse<CacheResponse>> {
        return this.client.get(`/${id}`)
    }

    async putById(id: string, accessList: CacheRequest): Promise<ApiResponse<void>> {
        return this.client.put(`/${id}`, accessList)
    }

    async deleteById(id: string): Promise<ApiResponse<void>> {
        return this.client.delete(`/${id}`)
    }

    async post(accessList: CacheRequest): Promise<ApiResponse<GenericCreateResponse>> {
        return this.client.post("", accessList)
    }
}
