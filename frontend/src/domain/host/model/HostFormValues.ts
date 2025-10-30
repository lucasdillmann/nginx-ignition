import { CertificateResponse } from "../../certificate/model/CertificateResponse"
import { HostBindingType, HostFeatureSet, HostRouteSettings, HostRouteSourceCode, HostRouteType } from "./HostRequest"
import IntegrationOptionResponse from "../../integration/model/IntegrationOptionResponse"
import AccessListResponse from "../../accesslist/model/AccessListResponse"
import IntegrationResponse from "../../integration/model/IntegrationResponse"

export interface HostFormBinding {
    type: HostBindingType
    ip: string
    port: number
    certificate?: CertificateResponse
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

export default interface HostFormValues {
    enabled: boolean
    defaultServer: boolean
    useGlobalBindings: boolean
    domainNames: string[]
    routes: HostFormRoute[]
    bindings: HostFormBinding[]
    featureSet: HostFeatureSet
    accessList?: AccessListResponse
}
