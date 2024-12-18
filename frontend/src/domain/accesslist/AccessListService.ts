import AccessListGateway from "./AccessListGateway"
import { requireNullablePayload, requireSuccessPayload, requireSuccessResponse } from "../../core/apiclient/ApiResponse"
import PageResponse from "../../core/pagination/PageResponse"
import AccessListRequest from "./model/AccessListRequest"
import AccessListResponse from "./model/AccessListResponse"

export default class AccessListService {
    private readonly gateway: AccessListGateway

    constructor() {
        this.gateway = new AccessListGateway()
    }

    async list(
        pageSize?: number,
        pageNumber?: number,
        searchTerms?: string,
    ): Promise<PageResponse<AccessListResponse>> {
        return this.gateway.getPage(pageSize, pageNumber, searchTerms).then(requireSuccessPayload)
    }

    async delete(id: string): Promise<void> {
        return this.gateway.deleteById(id).then(requireSuccessResponse)
    }

    async getById(id: string): Promise<AccessListResponse | undefined> {
        return this.gateway.getById(id).then(requireNullablePayload)
    }

    async updateById(id: string, user: AccessListRequest): Promise<void> {
        return this.gateway.putById(id, user).then(requireSuccessResponse)
    }

    async create(user: AccessListRequest): Promise<void> {
        return this.gateway.post(user).then(requireSuccessResponse)
    }
}
