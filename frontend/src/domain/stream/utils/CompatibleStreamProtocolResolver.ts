import { StreamProtocol } from "../model/StreamRequest"

export default class CompatibleStreamProtocolResolver {
    static resolve(protocol: StreamProtocol): StreamProtocol[] {
        return Object.values(StreamProtocol).filter(item => item == StreamProtocol.SOCKET || item == protocol)
    }
}
