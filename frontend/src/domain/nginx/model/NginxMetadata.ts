export enum NginxSupportType {
    NONE = "NONE",
    STATIC = "STATIC",
    DYNAMIC = "DYNAMIC",
}

export interface NginxAvailableSupport {
    runCode: NginxSupportType
    streams: NginxSupportType
    tlsSni: NginxSupportType
    stats: NginxSupportType
}

export interface NginxStatsConfig {
    enabled: boolean
    allHosts: boolean
}

export default interface NginxMetadata {
    version: string
    availableSupport: NginxAvailableSupport
    stats: NginxStatsConfig
}
