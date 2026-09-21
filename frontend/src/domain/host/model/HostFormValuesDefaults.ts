import HostFormValues from "./HostFormValues"
import { HostBindingType, HostRouteProtocol, HostRouteType } from "./HostRequest"

export function hostFormValuesDefaults(): HostFormValues {
    return {
        enabled: true,
        defaultServer: false,
        useGlobalBindings: true,
        domainNames: [""],
        vpns: [],
        bindings: [
            {
                ip: "0.0.0.0",
                port: 8080,
                type: HostBindingType.HTTP,
            },
        ],
        routes: [
            {
                priority: 0,
                enabled: true,
                type: HostRouteType.PROXY,
                protocol: HostRouteProtocol.HTTP_1_1,
                sourcePath: "/",
                targetUri: "",
                settings: {
                    keepOriginalDomainName: true,
                    proxySslServerName: true,
                    ignoreSslErrors: false,
                    includeForwardHeaders: true,
                    directoryListingEnabled: false,
                },
            },
        ],
        featureSet: {
            websocketsSupport: true,
            http2Support: true,
            redirectHttpToHttps: false,
            statsEnabled: true,
        },
    }
}
