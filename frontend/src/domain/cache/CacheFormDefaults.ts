import CacheRequest, { HttpMethod, UseStale } from "./model/CacheRequest"

export function cacheFormDefaults(): CacheRequest {
    return {
        name: "",
        allowedMethods: [HttpMethod.GET, HttpMethod.HEAD],
        minimumUsesBeforeCaching: 1,
        useStale: [UseStale.ERROR, UseStale.TIMEOUT, UseStale.UPDATING],
        backgroundUpdate: false,
        revalidate: true,
        cacheStatusResponseHeaderEnabled: true,
        ignoreUpstreamCacheHeaders: false,
        concurrencyLock: {
            enabled: false,
        },
        bypassRules: [],
        noCacheRules: [],
        fileExtensions: [],
        durations: [],
    }
}
