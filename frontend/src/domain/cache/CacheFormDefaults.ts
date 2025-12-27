import CacheRequest from "./model/CacheRequest"

export function cacheFormDefaults(): CacheRequest {
    return {
        name: "",
        storagePath: "",
        inactiveSeconds: 0,
        maximumSizeMb: 0,
        allowedMethods: [],
        minimumUsesBeforeCaching: 0,
        useStale: [],
        backgroundUpdate: false,
        concurrencyLock: {
            timeoutSeconds: 0,
            ageSeconds: 0,
            enabled: false,
        },
        revalidate: false,
        bypassRules: [],
        noCacheRules: [],
        durations: [],
    }
}
