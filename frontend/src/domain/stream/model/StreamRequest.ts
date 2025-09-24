export enum StreamProtocol {
    TCP = "TCP",
    UDP = "UDP",
    SOCKET = "SOCKET",
}

export enum StreamType {
    SIMPLE = "SIMPLE",
    SNI_ROUTER = "SNI_ROUTER",
}

export interface StreamAddress {
    protocol: StreamProtocol
    address: string
    port?: number
}

export interface StreamFeatureSet {
    useProxyProtocol: boolean
    socketKeepAlive: boolean
    tcpKeepAlive: boolean
    tcpNoDelay: boolean
    tcpDeferred: boolean
}

export interface StreamCircuitBreaker {
    maxFailures: number
    openSeconds: number
}

export interface StreamBackend {
    weight?: number
    target: StreamAddress
    circuitBreaker?: StreamCircuitBreaker
}

export interface StreamRoute {
    domainNames: string[]
    backends: StreamBackend[]
}

export default interface StreamRequest {
    enabled: boolean
    name: string
    type: StreamType
    featureSet: StreamFeatureSet
    defaultBackend: StreamBackend
    binding: StreamAddress
    routes?: StreamRoute[]
}
