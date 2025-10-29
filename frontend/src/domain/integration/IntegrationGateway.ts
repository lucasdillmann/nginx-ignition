import ApiClient from "../../core/apiclient/ApiClient"
import ApiResponse from "../../core/apiclient/ApiResponse"
import PageResponse from "../../core/pagination/PageResponse"
import IntegrationOptionResponse from "./model/IntegrationOptionResponse"
import IntegrationRequest from "./model/IntegrationRequest"
import GenericCreateResponse from "../../core/common/GenericCreateResponse"
import IntegrationResponse from "./model/IntegrationResponse"
import AvailableDriverResponse from "./model/AvailableDriverResponse"

export default class IntegrationGateway {
    private readonly client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/integrations")
    }

    async getById(id: string): Promise<ApiResponse<IntegrationResponse>> {
        return this.client.get(`/${id}`)
    }

    async getPage(
        pageSize?: number,
        pageNumber?: number,
        searchTerms?: string,
        enabledOnly?: boolean,
    ): Promise<ApiResponse<PageResponse<IntegrationResponse>>> {
        return this.client.get(undefined, undefined, { pageSize, pageNumber, searchTerms, enabledOnly })
    }

    async putById(id: string, user: IntegrationRequest): Promise<ApiResponse<void>> {
        return this.client.put(`/${id}`, user)
    }

    async post(user: IntegrationRequest): Promise<ApiResponse<GenericCreateResponse>> {
        return this.client.post("", user)
    }

    async delete(id: string): Promise<ApiResponse<void>> {
        return this.client.delete(`/${id}`)
    }

    async getAvailableDrivers(): Promise<ApiResponse<AvailableDriverResponse[]>> {
        return this.client.get(`/available-drivers`)
    }

    async getIntegrationOptions(
        id: string,
        pageSize?: number,
        pageNumber?: number,
        searchTerms?: string,
        tcpOnly?: boolean,
    ): Promise<ApiResponse<PageResponse<IntegrationOptionResponse>>> {
        return this.client.get(`/${id}/options`, undefined, { pageSize, pageNumber, searchTerms, tcpOnly })
    }

    async getIntegrationOptionById(
        integrationId: string,
        optionId: string,
    ): Promise<ApiResponse<IntegrationOptionResponse>> {
        return this.client.get(`/${integrationId}/options/${optionId}`)
    }
}
