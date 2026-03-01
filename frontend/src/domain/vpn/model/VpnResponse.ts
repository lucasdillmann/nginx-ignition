import EndpointSSLSupport from "./EndpointSSLSupport"
import VpnRequest from "./VpnRequest"

export default interface VpnResponse extends VpnRequest {
    id: string
    driverEndpointSslSupport: EndpointSSLSupport
}
