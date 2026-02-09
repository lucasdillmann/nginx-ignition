export default interface TrafficStatsResponse {
    serverZones: Record<string, ZoneData>
    filterZones: Record<string, Record<string, ZoneData>>
    upstreamZones: Record<string, UpstreamZoneData[]>
    hostName: string
    connections: ConnectionsData
}

export interface ZoneData {
    requestMsecs: TimeSeries
    responses: Responses
    requestCounter: number
    inBytes: number
    outBytes: number
    requestMsec: number
    requestMsecCounter: number
}

export interface UpstreamZoneData {
    server: string
    requestMsecs: TimeSeries
    responseMsecs: TimeSeries
    responses: UpstreamResponses
    requestMsecCounter: number
    responseMsec: number
    inBytes: number
    requestMsec: number
    responseMsecCounter: number
    requestCounter: number
    weight: number
    maxFails: number
    failTimeout: number
    outBytes: number
    backup: boolean
    down: boolean
}

export interface ConnectionsData {
    active: number
    reading: number
    writing: number
    waiting: number
    accepted: number
    handled: number
    requests: number
}

export interface Responses {
    "1xx": number
    "2xx": number
    "3xx": number
    "4xx": number
    "5xx": number
    miss: number
    bypass: number
    expired: number
    stale: number
    updating: number
    revalidated: number
    hit: number
    scarce: number
}

export interface UpstreamResponses {
    "1xx": number
    "2xx": number
    "3xx": number
    "4xx": number
    "5xx": number
}

export interface TimeSeries {
    times: number[]
    msecs: number[]
}
