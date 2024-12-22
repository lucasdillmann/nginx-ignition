import { AccessListCredentials, AccessListOutcome } from "./model/AccessListRequest"
import AccessListFormValues, { AccessListEntrySetFormValues } from "./model/AccessListFormValues"

const AccessListFormDefaults: AccessListFormValues = {
    name: "",
    realm: "",
    satisfyAll: true,
    defaultOutcome: AccessListOutcome.DENY,
    forwardAuthenticationHeader: false,
    credentials: [],
    entries: [],
}
export default AccessListFormDefaults

export const AccessListFormEntryDefaults: AccessListEntrySetFormValues = {
    priority: 0,
    outcome: AccessListOutcome.ALLOW,
    sourceAddresses: "",
}

export const AccessListFormCredentialsDefaults: AccessListCredentials = {
    username: "",
    password: "",
}
