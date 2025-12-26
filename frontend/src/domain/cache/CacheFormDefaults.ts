import { AccessListCredentials, AccessListOutcome } from "./model/CacheRequest"
import CacheFormValues, { AccessListEntrySetFormValues } from "./model/CacheFormValues"

export function cacheFormDefaults(): CacheFormValues {
    return {
        name: "",
        realm: "",
        satisfyAll: true,
        defaultOutcome: AccessListOutcome.DENY,
        forwardAuthenticationHeader: false,
        credentials: [],
        entries: [],
    }
}

export function accessListFormEntryDefaults(): AccessListEntrySetFormValues {
    return {
        priority: 0,
        outcome: AccessListOutcome.ALLOW,
        sourceAddresses: "",
    }
}

export function accessListFormCredentialsDefaults(): AccessListCredentials {
    return {
        username: "",
        password: "",
    }
}
