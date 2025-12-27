export enum HttpMethod {
    GET = "GET",
    HEAD = "HEAD",
    POST = "POST",
    PUT = "PUT",
    DELETE = "DELETE",
    PATCH = "PATCH",
    OPTIONS = "OPTIONS",
}

export enum UseStale {
    ERROR = "ERROR",
    TIMEOUT = "TIMEOUT",
    INVALID_HEADER = "INVALID_HEADER",
    UPDATING = "UPDATING",
    HTTP_500 = "HTTP_500",
    HTTP_502 = "HTTP_502",
    HTTP_503 = "HTTP_503",
    HTTP_504 = "HTTP_504",
    HTTP_403 = "HTTP_403",
    HTTP_404 = "HTTP_404",
    HTTP_429 = "HTTP_429",
}

export interface ConcurrencyLock {
    timeoutSeconds?: number
    ageSeconds?: number
    enabled: boolean
}

export interface Duration {
    statusCodes: number[]
    validTimeSeconds: number
}

export default interface CacheRequest {
    name: string
    storagePath?: string
    inactiveSeconds?: number
    maximumSizeMb?: number
    allowedMethods: HttpMethod[]
    minimumUsesBeforeCaching: number
    useStale: UseStale[]
    backgroundUpdate: boolean
    concurrencyLock: ConcurrencyLock
    revalidate: boolean
    cacheStatusResponseHeaderEnabled: boolean
    ignoreUpstreamCacheHeaders: boolean
    bypassRules: string[]
    noCacheRules: string[]
    fileExtensions: string[]
    durations: Duration[]
}
