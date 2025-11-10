import ApiClient from "../../core/apiclient/ApiClient"
import ApiResponse from "../../core/apiclient/ApiResponse"
import PageResponse from "../../core/pagination/PageResponse"
import VpnRequest from "./model/VpnRequest"
import GenericCreateResponse from "../../core/common/GenericCreateResponse"
import VpnResponse from "./model/VpnResponse"
import AvailableDriverResponse from "./model/AvailableDriverResponse"

export default class VpnGateway {
    private readonly client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/vpns")
    }

    async getById(id: string): Promise<ApiResponse<VpnResponse>> {
        return this.client.get(`/${id}`)
    }

    async getPage(
        pageSize?: number,
        pageNumber?: number,
        searchTerms?: string,
        enabledOnly?: boolean,
    ): Promise<ApiResponse<PageResponse<VpnResponse>>> {
        return this.client.get(undefined, undefined, { pageSize, pageNumber, searchTerms, enabledOnly })
    }

    async putById(id: string, user: VpnRequest): Promise<ApiResponse<void>> {
        return this.client.put(`/${id}`, user)
    }

    async post(user: VpnRequest): Promise<ApiResponse<GenericCreateResponse>> {
        return this.client.post("", user)
    }

    async delete(id: string): Promise<ApiResponse<void>> {
        return this.client.delete(`/${id}`)
    }

    async getAvailableDrivers(): Promise<ApiResponse<AvailableDriverResponse[]>> {
        return this.client.get("/available-drivers")
    }
}
