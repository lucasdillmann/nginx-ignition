export enum StreamProtocol {
    TCP = "TCP",
    UDP = "UDP",
    SOCKET = "SOCKET",
}

export interface StreamAddress {
    protocol: StreamProtocol
    address: string
    port?: number
}

export interface StreamFeatureSet {
    useProxyProtocol: boolean
    ssl: boolean
    tcpKeepAlive: boolean
    tcpNoDelay: boolean
    tcpDeferred: boolean
}

export default interface StreamRequest {
    enabled: boolean
    description: string
    featureSet: StreamFeatureSet
    binding: StreamAddress
    backend: StreamAddress
}
