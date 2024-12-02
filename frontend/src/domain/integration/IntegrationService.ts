import IntegrationGateway from "./IntegrationGateway"
import PageResponse from "../../core/pagination/PageResponse"
import { requireNullablePayload, requireSuccessPayload, requireSuccessResponse } from "../../core/apiclient/ApiResponse"
import { IntegrationResponse } from "./model/IntegrationResponse"
import { IntegrationConfigurationResponse } from "./model/IntegrationConfigurationResponse"
import { IntegrationConfigurationRequest } from "./model/IntegrationConfigurationRequest"
import { IntegrationOptionResponse } from "./model/IntegrationOptionResponse"

export default class IntegrationService {
    private readonly gateway: IntegrationGateway

    constructor() {
        this.gateway = new IntegrationGateway()
    }

    async getAll(enabledOnly: boolean = false): Promise<IntegrationResponse[]> {
        return this.gateway
            .getIntegrations()
            .then(requireSuccessPayload)
            .then(items => (enabledOnly ? items.filter(item => item.enabled) : items))
    }

    async getConfiguration(id: string): Promise<IntegrationConfigurationResponse> {
        return this.gateway.getIntegrationConfiguration(id).then(requireSuccessPayload)
    }

    async setConfiguration(id: string, configuration: IntegrationConfigurationRequest): Promise<void> {
        return this.gateway.putIntegrationConfiguration(id, configuration).then(requireSuccessResponse)
    }

    async getOptions(
        id: string,
        pageSize?: number,
        pageNumber?: number,
    ): Promise<PageResponse<IntegrationOptionResponse>> {
        return this.gateway.getIntegrationOptions(id, pageSize, pageNumber).then(requireSuccessPayload)
    }

    async getOptionById(integrationId: string, optionId: string): Promise<IntegrationOptionResponse | undefined> {
        return this.gateway.getIntegrationOptionById(integrationId, optionId).then(requireNullablePayload)
    }
}
