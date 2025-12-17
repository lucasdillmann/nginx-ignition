import HostRequest, { HostBinding } from "./HostRequest"

export default interface HostResponse extends HostRequest {
    id: string
    globalBindings?: HostBinding[]
}
