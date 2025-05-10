import ApiClient from "../../core/apiclient/ApiClient"
import ApiResponse from "../../core/apiclient/ApiResponse"
import PageResponse from "../../core/pagination/PageResponse"
import StreamResponse from "./model/StreamResponse"
import StreamRequest from "./model/StreamRequest"
import GenericCreateResponse from "../../core/common/GenericCreateResponse"

export default class StreamGateway {
    private readonly client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/streams")
    }

    async getPage(
        pageSize?: number,
        pageNumber?: number,
        searchTerms?: string,
    ): Promise<ApiResponse<PageResponse<StreamResponse>>> {
        return this.client.get(undefined, undefined, { pageSize, pageNumber, searchTerms })
    }

    async getById(id: string): Promise<ApiResponse<StreamResponse>> {
        return this.client.get(`/${id}`)
    }

    async putById(id: string, stream: StreamRequest): Promise<ApiResponse<void>> {
        return this.client.put(`/${id}`, stream)
    }

    async deleteById(id: string): Promise<ApiResponse<void>> {
        return this.client.delete(`/${id}`)
    }

    async post(stream: StreamRequest): Promise<ApiResponse<GenericCreateResponse>> {
        return this.client.post("", stream)
    }

    async toggleEnabled(id: string): Promise<ApiResponse<void>> {
        return this.client.post(`/${id}/toggle-enabled`)
    }
}
