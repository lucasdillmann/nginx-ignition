import ConfigurationGateway from "./ConfigurationGateway"
import { requireSuccessPayload } from "../../core/apiclient/ApiResponse"
import Configuration from "./Configuration"

const FALLBACK_RESPONSE: Configuration = {
    codeEditor: {
        apiKey: undefined,
    },
}

export default class ConfigurationService {
    private readonly gateway: ConfigurationGateway

    constructor() {
        this.gateway = new ConfigurationGateway()
    }

    public async get(): Promise<Configuration> {
        return this.gateway
            .get()
            .then(requireSuccessPayload)
            .catch(() => FALLBACK_RESPONSE)
    }
}
