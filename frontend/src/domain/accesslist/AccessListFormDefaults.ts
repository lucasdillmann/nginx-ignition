import { AccessListCredentials, AccessListOutcome } from "./model/AccessListRequest"
import AccessListFormValues, { AccessListEntrySetFormValues } from "./model/AccessListFormValues"

export function accessListFormDefaults(): AccessListFormValues {
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
