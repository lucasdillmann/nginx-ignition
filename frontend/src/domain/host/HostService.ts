import HostGateway from "./HostGateway";
import PageResponse from "../../core/pagination/PageResponse";
import HostResponse from "./model/HostResponse";
import {requireSuccessPayload, requireSuccessResponse} from "../../core/apiclient/ApiResponse";

export default class HostService {
    private readonly gateway: HostGateway

    constructor() {
        this.gateway = new HostGateway()
    }

    list(pageSize?: number, pageNumber?: number): Promise<PageResponse<HostResponse>> {
        return this.gateway.getPage(pageSize, pageNumber).then(requireSuccessPayload)
    }

    delete(id: string): Promise<void> {
        return this.gateway.delete(id).then(requireSuccessResponse)
    }

    toggleStatus(id: string): Promise<void> {
        return this.gateway.toggleStatus(id).then(requireSuccessResponse)
    }
}
