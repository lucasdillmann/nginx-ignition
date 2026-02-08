import HostFormValues from "./HostFormValues"
import { HostBindingType, HostRouteType } from "./HostRequest"

export function hostFormValuesDefaults(): HostFormValues {
    return {
        enabled: true,
        defaultServer: false,
        useGlobalBindings: true,
        statsEnabled: false,
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
                sourcePath: "/",
                targetUri: "",
                settings: {
                    keepOriginalDomainName: true,
                    proxySslServerName: true,
                    includeForwardHeaders: true,
                    directoryListingEnabled: false,
                },
            },
        ],
        featureSet: {
            websocketsSupport: true,
            http2Support: true,
            redirectHttpToHttps: false,
        },
    }
}
