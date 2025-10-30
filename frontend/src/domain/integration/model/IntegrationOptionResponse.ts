export enum IntegrationOptionProtocol {
    TCP = "TCP",
    UDP = "UDP",
}

export default interface IntegrationOptionResponse {
    id: string
    name: string
    port: number
    protocol: IntegrationOptionProtocol
}
