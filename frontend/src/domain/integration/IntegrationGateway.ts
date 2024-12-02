import ApiClient from "../../core/apiclient/ApiClient"
import ApiResponse from "../../core/apiclient/ApiResponse"
import PageResponse from "../../core/pagination/PageResponse"
import { IntegrationResponse } from "./model/IntegrationResponse"
import { IntegrationOptionResponse } from "./model/IntegrationOptionResponse"
import { IntegrationConfigurationResponse } from "./model/IntegrationConfigurationResponse"
import { IntegrationConfigurationRequest } from "./model/IntegrationConfigurationRequest"

export default class IntegrationGateway {
    private readonly client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/integrations")
    }

    async getIntegrations(): Promise<ApiResponse<IntegrationResponse[]>> {
        return this.client.get()
    }

    async getIntegrationOptions(
        id: string,
        pageSize?: number,
        pageNumber?: number,
    ): Promise<ApiResponse<PageResponse<IntegrationOptionResponse>>> {
        return this.client.get(`/${id}/options`, undefined, { pageSize, pageNumber })
    }

    async getIntegrationOptionById(
        integrationId: string,
        optionId: string,
    ): Promise<ApiResponse<IntegrationOptionResponse>> {
        return this.client.get(`/${integrationId}/options/${optionId}`)
    }

    async getIntegrationConfiguration(id: string): Promise<ApiResponse<IntegrationConfigurationResponse>> {
        return this.client.get(`/${id}/configuration`)
    }

    async putIntegrationConfiguration(
        id: string,
        payload: IntegrationConfigurationRequest,
    ): Promise<ApiResponse<void>> {
        return this.client.put(`/${id}/configuration`, payload)
    }
}
