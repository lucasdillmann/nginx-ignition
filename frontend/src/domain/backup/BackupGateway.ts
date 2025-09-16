import ApiClient from "../../core/apiclient/ApiClient"
import ApiResponse from "../../core/apiclient/ApiResponse"

export default class BackupGateway {
    private readonly client: ApiClient

    constructor() {
        this.client = new ApiClient("/api/backup")
    }

    async download(): Promise<ApiResponse<any>> {
        return this.client.get("", undefined, undefined, true)
    }
}
