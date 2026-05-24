import { UpstreamResponses, ZoneData } from "../model/TrafficStatsResponse"

export const STATUS_CODES = ["1xx", "2xx", "3xx", "4xx", "5xx"] as const

const CHART_PALETTE_SIZE = 8

function readThemeCssVar(name: string): string {
    return getComputedStyle(document.documentElement).getPropertyValue(name).trim()
}

export function getStatusChartColors(): { domain: string[]; range: string[] } {
    const domain = [...STATUS_CODES]
    const range = STATUS_CODES.map(status => readThemeCssVar(`--nginxIgnition-chartStatus${status}`))
    return { domain, range }
}

export function getChartPalette(): string[] {
    return Array.from({ length: CHART_PALETTE_SIZE }, (_, index) =>
        readThemeCssVar(`--nginxIgnition-chartPalette-${index + 1}`),
    )
}

export function getChartPrimaryColor(): string {
    return readThemeCssVar("--nginxIgnition-chartPrimary")
}

export function getChartColorScale(categories: string[]): { domain: string[]; range: string[] } {
    const palette = getChartPalette()
    const domain = [...categories]
    const range = domain.map((_, index) => palette[index % palette.length])
    return { domain, range }
}

export interface StatusDataItem {
    status: string
    count: number
}

export function buildStatusDistributionData(responses: UpstreamResponses): StatusDataItem[] {
    return [
        { status: "1xx", count: responses["1xx"] },
        { status: "2xx", count: responses["2xx"] },
        { status: "3xx", count: responses["3xx"] },
        { status: "4xx", count: responses["4xx"] },
        { status: "5xx", count: responses["5xx"] },
    ].filter(item => item.count > 0)
}

export interface TrafficDataItem {
    name: string
    requests: number
    inBytes: number
    outBytes: number
}

export function buildTrafficByDomainData(
    serverZones: Record<string, ZoneData> | null,
    topN: number = 10,
): TrafficDataItem[] {
    if (!serverZones) return []

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

export function aggregateResponses(zones: Record<string, ZoneData> | null): ZoneData["responses"] {
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

    if (!zones) return result

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
    if (!timeSeries?.times || timeSeries.times.length === 0) {
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

export function buildCountryCodeData(countryCodes: Record<string, ZoneData>) {
    if (!countryCodes) return []

    return Object.entries(countryCodes)
        .map(([country, data]) => ({
            country,
            value: data.requestCounter,
        }))
        .filter(item => item.value > 0)
        .sort((a, b) => b.value - a.value)
}

export function buildCityData(cities: Record<string, ZoneData>) {
    if (!cities) return []

    return Object.entries(cities)
        .map(([city, data]) => ({
            city,
            value: data.requestCounter,
        }))
        .filter(item => item.value > 0)
        .sort((a, b) => b.value - a.value)
}
