import { ServerCertificateResponse } from "../../certificate/server/model/ServerCertificateResponse"
import { HostBindingType, HostFeatureSet, HostRouteSettings, HostRouteSourceCode, HostRouteType } from "./HostRequest"
import IntegrationOptionResponse from "../../integration/model/IntegrationOptionResponse"
import AccessListResponse from "../../accesslist/model/AccessListResponse"
import IntegrationResponse from "../../integration/model/IntegrationResponse"
import VpnResponse from "../../vpn/model/VpnResponse"

export interface HostFormBinding {
    type: HostBindingType
    ip: string
    port: number
    certificate?: ServerCertificateResponse
}

export interface HostFormStaticResponse {
    statusCode: number
    payload?: string
    headers?: string
}

export interface HostFormRoute {
    priority: number
    enabled: boolean
    type: HostRouteType
    sourcePath: string
    settings: HostRouteSettings
    targetUri?: string
    response?: HostFormStaticResponse
    integration?: HostFormRouteIntegration
    accessList?: AccessListResponse
    redirectCode?: number
    sourceCode?: HostRouteSourceCode
}

export interface HostFormRouteIntegration {
    integration: IntegrationResponse
    option: IntegrationOptionResponse
}

export interface HostFormVpn {
    vpn: VpnResponse
    name: string
    host?: string
}

export default interface HostFormValues {
    enabled: boolean
    defaultServer: boolean
    useGlobalBindings: boolean
    domainNames: string[]
    routes: HostFormRoute[]
    bindings: HostFormBinding[]
    vpns: HostFormVpn[]
    featureSet: HostFeatureSet
    accessList?: AccessListResponse
}
