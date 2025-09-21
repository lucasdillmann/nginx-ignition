import StreamRequest, {
    StreamBackend,
    StreamCircuitBreaker,
    StreamProtocol,
    StreamRoute,
    StreamType,
} from "./model/StreamRequest"

export const StreamCircuitBreakerDefaults: StreamCircuitBreaker = Object.freeze({
    maxFailures: 5,
    openSeconds: 30,
})

export const StreamBackendDefault: StreamBackend = Object.freeze({
    target: {
        protocol: StreamProtocol.TCP,
        address: "",
        port: 8080,
    },
})

export const StreamRouteDefaults: StreamRoute = Object.freeze({
    domainNames: [""],
    backends: [StreamBackendDefault],
})

const StreamFormDefaults: StreamRequest = Object.freeze({
    name: "",
    enabled: true,
    type: StreamType.SIMPLE,
    featureSet: {
        useProxyProtocol: false,
        socketKeepAlive: true,
        tcpKeepAlive: false,
        tcpNoDelay: false,
        tcpDeferred: false,
    },
    binding: {
        protocol: StreamProtocol.TCP,
        address: "0.0.0.0",
        port: 8080,
    },
    defaultBackend: StreamBackendDefault,
})

export default StreamFormDefaults
