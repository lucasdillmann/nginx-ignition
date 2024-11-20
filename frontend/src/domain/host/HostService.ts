import HostGateway from "./HostGateway";
import PageResponse from "../../core/pagination/PageResponse";
import HostResponse from "./model/HostResponse";
import {requireSuccessPayload, requireSuccessResponse} from "../../core/apiclient/ApiResponse";

export default class HostService {
    private readonly gateway: HostGateway

    constructor() {
        this.gateway = new HostGateway()
    }

    async list(pageSize?: number, pageNumber?: number): Promise<PageResponse<HostResponse>> {
        return this.gateway.getPage(pageSize, pageNumber).then(requireSuccessPayload)
    }

    async delete(id: string): Promise<void> {
        return this.gateway.delete(id).then(requireSuccessResponse)
    }

    async toggleEnabled(id: string): Promise<void> {
        const response = await this.gateway.toggleEnabled(id);
        return requireSuccessResponse(response);
    }
}
