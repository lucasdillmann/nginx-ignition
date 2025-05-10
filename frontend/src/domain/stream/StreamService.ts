import StreamGateway from "./StreamGateway"
import { requireNullablePayload, requireSuccessPayload, requireSuccessResponse } from "../../core/apiclient/ApiResponse"
import PageResponse from "../../core/pagination/PageResponse"
import StreamRequest from "./model/StreamRequest"
import StreamResponse from "./model/StreamResponse"
import GenericCreateResponse from "../../core/common/GenericCreateResponse"

export default class StreamService {
    private readonly gateway: StreamGateway

    constructor() {
        this.gateway = new StreamGateway()
    }

    async list(pageSize?: number, pageNumber?: number, searchTerms?: string): Promise<PageResponse<StreamResponse>> {
        return this.gateway.getPage(pageSize, pageNumber, searchTerms).then(requireSuccessPayload)
    }

    async delete(id: string): Promise<void> {
        return this.gateway.deleteById(id).then(requireSuccessResponse)
    }

    async getById(id: string): Promise<StreamResponse | undefined> {
        return this.gateway.getById(id).then(requireNullablePayload)
    }

    async updateById(id: string, user: StreamRequest): Promise<void> {
        return this.gateway.putById(id, user).then(requireSuccessResponse)
    }

    async create(user: StreamRequest): Promise<GenericCreateResponse> {
        return this.gateway.post(user).then(requireSuccessPayload)
    }

    async toggleEnabled(id: string): Promise<void> {
        return this.gateway.toggleEnabled(id).then(requireSuccessResponse)
    }
}
