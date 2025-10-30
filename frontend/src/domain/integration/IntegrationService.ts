import IntegrationGateway from "./IntegrationGateway"
import PageResponse from "../../core/pagination/PageResponse"
import { requireNullablePayload, requireSuccessPayload, requireSuccessResponse } from "../../core/apiclient/ApiResponse"
import GenericCreateResponse from "../../core/common/GenericCreateResponse"
import IntegrationRequest from "./model/IntegrationRequest"
import IntegrationResponse from "./model/IntegrationResponse"
import IntegrationOptionResponse from "./model/IntegrationOptionResponse"
import AvailableDriverResponse from "./model/AvailableDriverResponse"

export default class IntegrationService {
    private readonly gateway: IntegrationGateway

    constructor() {
        this.gateway = new IntegrationGateway()
    }

    async getById(id: string): Promise<IntegrationResponse | undefined> {
        return this.gateway.getById(id).then(requireNullablePayload)
    }

    async list(
        pageSize?: number,
        pageNumber?: number,
        searchTerms?: string,
        enabledOnly?: boolean,
    ): Promise<PageResponse<IntegrationResponse>> {
        return this.gateway.getPage(pageSize, pageNumber, searchTerms, enabledOnly).then(requireSuccessPayload)
    }

    async delete(id: string): Promise<void> {
        return this.gateway.delete(id).then(requireSuccessResponse)
    }

    async updateById(id: string, integration: IntegrationRequest): Promise<void> {
        return this.gateway.putById(id, integration).then(requireSuccessResponse)
    }

    async create(integration: IntegrationRequest): Promise<GenericCreateResponse> {
        return this.gateway.post(integration).then(requireSuccessPayload)
    }

    async availableDrivers(): Promise<AvailableDriverResponse[]> {
        return this.gateway.getAvailableDrivers().then(requireSuccessPayload)
    }

    async getOptions(
        id: string,
        pageSize?: number,
        pageNumber?: number,
        searchTerms?: string,
        tcpOnly?: boolean,
    ): Promise<PageResponse<IntegrationOptionResponse>> {
        return this.gateway
            .getIntegrationOptions(id, pageSize, pageNumber, searchTerms, tcpOnly)
            .then(requireSuccessPayload)
    }

    async getOptionById(integrationId: string, optionId: string): Promise<IntegrationOptionResponse | undefined> {
        return this.gateway.getIntegrationOptionById(integrationId, optionId).then(requireNullablePayload)
    }
}
