export enum HostBindingType {
    HTTP = "HTTP",
    HTTPS = "HTTPS",
}

export enum HostRouteType {
    PROXY = "PROXY",
    REDIRECT = "REDIRECT",
    STATIC_RESPONSE = "STATIC_RESPONSE",
    INTEGRATION = "INTEGRATION",
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
    certificateId?: string
}

export interface HostRouteStaticResponse {
    statusCode: number
    payload?: string
    headers?: Record<string, string>
}

export interface HostRoute {
    priority: number
    type: HostRouteType
    sourcePath: string
    targetUri?: string
    customSettings?: string
    response?: HostRouteStaticResponse
    integration?: HostRouteIntegration
}

export interface HostRouteIntegration {
    integrationId: string
    optionId: string
}

export default interface HostRequest {
    enabled: boolean
    defaultServer: boolean
    domainNames: string[]
    routes: HostRoute[]
    bindings: HostBinding[]
    featureSet: HostFeatureSet
}
