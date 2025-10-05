export enum NginxSupportType {
    NONE = "NONE",
    STATIC = "STATIC",
    DYNAMIC = "DYNAMIC",
}

export interface NginxAvailableSupport {
    runCode: NginxSupportType
    streams: NginxSupportType
    tlsSni: NginxSupportType
}

export default interface NginxMetadata {
    version: string
    availableSupport: NginxAvailableSupport
}
