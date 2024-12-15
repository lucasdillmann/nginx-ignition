import HostFormValues from "./HostFormValues"
import { HostBindingType, HostRouteType } from "./HostRequest"

const HostFormValuesDefaults: HostFormValues = {
    enabled: true,
    defaultServer: false,
    useGlobalBindings: true,
    domainNames: [""],
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
            type: HostRouteType.PROXY,
            sourcePath: "/",
            targetUri: "",
            settings: {
                forwardQueryParams: true,
                keepOriginalDomainName: true,
                proxySslServerName: true,
                includeForwardHeaders: true,
            },
        },
    ],
    featureSet: {
        websocketsSupport: true,
        http2Support: true,
        redirectHttpToHttps: false,
    },
}

export default HostFormValuesDefaults
