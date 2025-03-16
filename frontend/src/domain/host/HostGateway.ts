import ApiClient from "../../core/apiclient/ApiClient"
import ApiResponse from "../../core/apiclient/ApiResponse"
import PageResponse from "../../core/pagination/PageResponse"
import HostResponse from "./model/HostResponse"
import HostRequest from "./model/HostRequest"
import GenericCreateResponse from "../../core/common/GenericCreateResponse"

export default class HostGateway {
    private readonly client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/hosts")
    }

    async getById(id: string): Promise<ApiResponse<HostResponse>> {
        return this.client.get(`/${id}`)
    }

    async getPage(
        pageSize?: number,
        pageNumber?: number,
        searchTerms?: string,
    ): Promise<ApiResponse<PageResponse<HostResponse>>> {
        return this.client.get(undefined, undefined, { pageSize, pageNumber, searchTerms })
    }

    async delete(id: string): Promise<ApiResponse<void>> {
        return this.client.delete(`/${id}`)
    }

    async toggleEnabled(id: string): Promise<ApiResponse<void>> {
        return this.client.post(`/${id}/toggle-enabled`)
    }

    async getLogs(id: string, type: string, lines: number): Promise<ApiResponse<string[]>> {
        return this.client.get(`/${id}/logs/${type}`, undefined, { lines })
    }

    async putById(id: string, user: HostRequest): Promise<ApiResponse<void>> {
        return this.client.put(`/${id}`, user)
    }

    async post(user: HostRequest): Promise<ApiResponse<GenericCreateResponse>> {
        return this.client.post("", user)
    }
}
