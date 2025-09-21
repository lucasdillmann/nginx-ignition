import StreamRequest, { StreamProtocol, StreamType } from "./model/StreamRequest"

const StreamFormDefaults: StreamRequest = {
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
    defaultBackend: {
        target: {
            protocol: StreamProtocol.TCP,
            address: "",
            port: 8080,
        },
    },
}

export default StreamFormDefaults
