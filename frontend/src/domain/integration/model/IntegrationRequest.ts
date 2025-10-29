export enum IntegrationDriver {
    DOCKER = "DOCKER",
    TRUENAS = "TRUENAS",
}

export default interface IntegrationRequest {
    id: string
    name: string
    driver: IntegrationDriver
    enabled: boolean
    parameters: object
}
