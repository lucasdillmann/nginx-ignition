export enum HostBindingType {
    HTTP = "HTTP",
    HTTPS = "HTTPS",
}

export enum HostRouteType {
    PROXY = "PROXY",
    REDIRECT = "REDIRECT",
    STATIC_RESPONSE = "STATIC_RESPONSE",
    INTEGRATION = "INTEGRATION",
    EXECUTE_CODE = "EXECUTE_CODE",
    STATIC_FILES = "STATIC_FILES",
}

export enum HostRouteSourceCodeLanguage {
    JAVASCRIPT = "JAVASCRIPT",
    LUA = "LUA",
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

export interface HostRouteSettings {
    includeForwardHeaders: boolean
    proxySslServerName: boolean
    keepOriginalDomainName: boolean
    directoryListingEnabled: boolean
    custom?: string
}

export interface HostRouteSourceCode {
    language: HostRouteSourceCodeLanguage
    code: string
    mainFunction?: string
}

export interface HostRoute {
    priority: number
    enabled: boolean
    type: HostRouteType
    sourcePath: string
    settings: HostRouteSettings
    targetUri?: string
    response?: HostRouteStaticResponse
    integration?: HostRouteIntegration
    accessListId?: string
    redirectCode?: number
    sourceCode?: HostRouteSourceCode
}

export interface HostRouteIntegration {
    integrationId: string
    optionId: string
}

export default interface HostRequest {
    enabled: boolean
    defaultServer: boolean
    useGlobalBindings: boolean
    domainNames?: string[]
    routes: HostRoute[]
    bindings: HostBinding[]
    featureSet: HostFeatureSet
    accessListId?: string
}
