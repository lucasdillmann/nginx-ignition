import IntegrationRequest from "./IntegrationRequest"

export function integrationRequestDefaults(): IntegrationRequest {
    return {
        enabled: true,
        name: "",
        driver: "",
        parameters: {},
    }
}
