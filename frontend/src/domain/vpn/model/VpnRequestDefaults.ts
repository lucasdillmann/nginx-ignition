import VpnRequest from "./VpnRequest"

export function vpnRequestDefaults(): VpnRequest {
    return {
        enabled: true,
        name: "",
        driver: "",
        parameters: {},
    }
}
