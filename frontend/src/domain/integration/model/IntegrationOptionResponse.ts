export enum IntegrationOptionProtocol {
    TCP = "TCP",
    UDP = "UDP",
}

export interface IntegrationOptionResponse {
    id: string
    name: string
    port: number
    protocol: IntegrationOptionProtocol
}
