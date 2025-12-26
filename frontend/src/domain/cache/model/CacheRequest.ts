export enum AccessListOutcome {
    ALLOW = "ALLOW",
    DENY = "DENY",
}

export interface AccessListEntrySet {
    priority: number
    outcome: AccessListOutcome
    sourceAddresses: string[]
}

export interface AccessListCredentials {
    username: string
    password: string
}

export default interface CacheRequest {
    name: string
    realm?: string
    satisfyAll: boolean
    defaultOutcome: AccessListOutcome
    entries?: AccessListEntrySet[]
    credentials?: AccessListCredentials[]
    forwardAuthenticationHeader: boolean
}
