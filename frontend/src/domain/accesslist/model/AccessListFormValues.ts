import { AccessListCredentials, AccessListOutcome } from "./AccessListRequest"

export interface AccessListEntrySetFormValues {
    priority: number
    outcome: AccessListOutcome
    sourceAddresses: string
}

export default interface AccessListFormValues {
    name: string
    realm?: string
    satisfyAll: boolean
    defaultOutcome: AccessListOutcome
    entries: AccessListEntrySetFormValues[]
    credentials: AccessListCredentials[]
    forwardAuthenticationHeader: boolean
}
