import { ZoneData } from "../model/TrafficStatsResponse"

export const STATUS_COLORS = {
    "1xx": "#1677ff",
    "2xx": "#52c41a",
    "3xx": "#faad14",
    "4xx": "#ff7a45",
    "5xx": "#f5222d",
}

export interface StatusDataItem {
    status: string
    count: number
    color: string
}

export function buildStatusDistributionData(responses: ZoneData["responses"]): StatusDataItem[] {
    return [
        { status: "1xx", count: responses["1xx"], color: STATUS_COLORS["1xx"] },
        { status: "2xx", count: responses["2xx"], color: STATUS_COLORS["2xx"] },
        { status: "3xx", count: responses["3xx"], color: STATUS_COLORS["3xx"] },
        { status: "4xx", count: responses["4xx"], color: STATUS_COLORS["4xx"] },
        { status: "5xx", count: responses["5xx"], color: STATUS_COLORS["5xx"] },
    ].filter(item => item.count > 0)
}

export interface TrafficDataItem {
    name: string
    requests: number
    inBytes: number
    outBytes: number
}

export function buildTrafficByDomainData(serverZones: Record<string, ZoneData>, topN: number = 10): TrafficDataItem[] {
    return Object.entries(serverZones)
        .filter(([name]) => name !== "*")
        .map(([name, data]) => ({
            name,
            requests: data.requestCounter,
            inBytes: data.inBytes,
            outBytes: data.outBytes,
        }))
        .sort((a, b) => b.requests - a.requests)
        .slice(0, topN)
}

export function aggregateResponses(zones: Record<string, ZoneData>): ZoneData["responses"] {
    const result = {
        "1xx": 0,
        "2xx": 0,
        "3xx": 0,
        "4xx": 0,
        "5xx": 0,
        miss: 0,
        bypass: 0,
        expired: 0,
        stale: 0,
        updating: 0,
        revalidated: 0,
        hit: 0,
        scarce: 0,
    }

    Object.values(zones).forEach(zone => {
        result["1xx"] += zone.responses["1xx"]
        result["2xx"] += zone.responses["2xx"]
        result["3xx"] += zone.responses["3xx"]
        result["4xx"] += zone.responses["4xx"]
        result["5xx"] += zone.responses["5xx"]
        result.miss += zone.responses.miss
        result.bypass += zone.responses.bypass
        result.expired += zone.responses.expired
        result.stale += zone.responses.stale
        result.updating += zone.responses.updating
        result.revalidated += zone.responses.revalidated
        result.hit += zone.responses.hit
        result.scarce += zone.responses.scarce
    })

    return result
}

export function buildResponseTimeData(timeSeries: { times: number[]; msecs: number[] }) {
    if (!timeSeries || !timeSeries.times || timeSeries.times.length === 0) {
        return []
    }

    return timeSeries.times
        .map((time, index) => ({
            time: new Date(time).toLocaleTimeString(),
            timestamp: time,
            value: timeSeries.msecs[index],
        }))
        .filter(item => item.value > 0)
        .sort((a, b) => a.timestamp - b.timestamp)
}

export function buildUserAgentData(userAgents: Record<string, ZoneData>) {
    if (!userAgents) return []

    return Object.entries(userAgents)
        .map(([agent, data]) => ({
            type: agent,
            value: data.requestCounter,
        }))
        .filter(item => item.value > 0)
        .sort((a, b) => b.value - a.value)
}
