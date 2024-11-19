import ApiClient from "../../core/apiclient/ApiClient";
import ApiResponse from "../../core/apiclient/ApiResponse";
import PageResponse from "../../core/pagination/PageResponse";
import HostResponse from "./model/HostResponse";

export default class HostGateway {
    private readonly client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/hosts")
    }

    getPage(pageSize?: number, pageNumber?: number): Promise<ApiResponse<PageResponse<HostResponse>>> {
        return this.client.get(undefined, undefined, { pageSize, pageNumber })
    }

    delete(id: string): Promise<ApiResponse<void>> {
        return this.client.delete(`/${id}`)
    }

    toggleStatus(id: string): Promise<ApiResponse<void>> {
        return this.client.post(`/${id}/toggle-status`)
    }
}
