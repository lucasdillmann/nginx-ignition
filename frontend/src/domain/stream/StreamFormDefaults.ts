import StreamRequest, {
    StreamBackend,
    StreamCircuitBreaker,
    StreamProtocol,
    StreamRoute,
    StreamType,
} from "./model/StreamRequest"

export function streamCircuitBreakerDefaults(): StreamCircuitBreaker {
    return {
        maxFailures: 5,
        openSeconds: 30,
    }
}

export function streamBackendDefaults(): StreamBackend {
    return {
        target: {
            protocol: StreamProtocol.TCP,
            address: "",
            port: 8080,
        },
    }
}
export function streamRouteDefaults(): StreamRoute {
    return {
        domainNames: [""],
        backends: [streamBackendDefaults()],
    }
}

export function streamFormDefaults(): StreamRequest {
    return {
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
        defaultBackend: streamBackendDefaults(),
    }
}
