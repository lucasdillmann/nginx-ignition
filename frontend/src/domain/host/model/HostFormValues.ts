import { CertificateResponse } from "../../certificate/model/CertificateResponse"
import { HostBindingType, HostFeatureSet, HostRouteType } from "./HostRequest"

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
    type: HostRouteType
    sourcePath: string
    targetUri?: string
    customSettings?: string
    response?: HostFormStaticResponse
}

export default interface HostFormValues {
    enabled: boolean
    defaultServer: boolean
    domainNames: string[]
    routes: HostFormRoute[]
    bindings: HostFormBinding[]
    featureSet: HostFeatureSet
}
