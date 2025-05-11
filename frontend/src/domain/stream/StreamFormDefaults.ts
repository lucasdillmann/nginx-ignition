import StreamRequest, { StreamProtocol } from "./model/StreamRequest"

const StreamFormDefaults: StreamRequest = {
    name: "",
    enabled: true,
    featureSet: {
        useProxyProtocol: false,
        socketKeepAlive: false,
        tcpKeepAlive: false,
        tcpNoDelay: false,
        tcpDeferred: false,
    },
    binding: {
        protocol: StreamProtocol.TCP,
        address: "0.0.0.0",
        port: 8080,
    },
    backend: {
        protocol: StreamProtocol.TCP,
        address: "",
        port: 8080,
    },
}

export default StreamFormDefaults
