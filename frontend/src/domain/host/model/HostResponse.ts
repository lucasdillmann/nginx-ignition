export enum HostBindingType {
    HTTP = "HTTP",
    HTTPS = "HTTPS",
}

export enum HostRouteType {
    PROXY = "PROXY",
    REDIRECT = "REDIRECT",
    STATIC_RESPONSE = "STATIC_RESPONSE",
}

export interface HostFeatureSet {
    websocketsSupport: boolean
    http2Support: boolean
    redirectHttpToHttps: boolean
}

export interface HostBinding {
    type: HostBindingType
    ip: string
    port: number
    certificateId: string
}

export interface HostRouteStaticResponse {
    statusCode: number
    payload?: string
    headers?: { [key: string]: string[] }
}

export interface HostRoute {
    priority: number
    type: HostRouteType
    sourcePath: string
    targetUri?: string
    customSettings?: string
    response?: HostRouteStaticResponse
}

export default interface HostResponse {
    id: string
    default: boolean
    enabled: boolean
    domainNames: string[]
    routes: HostRoute[]
    binding: HostBinding[]
    featureSet: HostFeatureSet
}