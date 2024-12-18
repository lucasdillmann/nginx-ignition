import ApiClient from "../../core/apiclient/ApiClient"
import ApiResponse from "../../core/apiclient/ApiResponse"
import PageResponse from "../../core/pagination/PageResponse"
import AccessListResponse from "./model/AccessListResponse"
import AccessListRequest from "./model/AccessListRequest"

export default class AccessListGateway {
    private readonly client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/access-lists")
    }

    async getPage(
        pageSize?: number,
        pageNumber?: number,
        searchTerms?: string,
    ): Promise<ApiResponse<PageResponse<AccessListResponse>>> {
        return this.client.get(undefined, undefined, { pageSize, pageNumber, searchTerms })
    }

    async getById(id: string): Promise<ApiResponse<AccessListResponse>> {
        return this.client.get(`/${id}`)
    }

    async putById(id: string, accessList: AccessListRequest): Promise<ApiResponse<void>> {
        return this.client.put(`/${id}`, accessList)
    }

    async deleteById(id: string): Promise<ApiResponse<void>> {
        return this.client.delete(`/${id}`)
    }

    async post(accessList: AccessListRequest): Promise<ApiResponse<void>> {
        return this.client.post("", accessList)
    }
}
