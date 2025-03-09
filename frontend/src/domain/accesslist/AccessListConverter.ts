import AccessListResponse from "./model/AccessListResponse"
import AccessListFormValues, { AccessListEntrySetFormValues } from "./model/AccessListFormValues"
import AccessListRequest, { AccessListEntrySet } from "./model/AccessListRequest"

class AccessListConverter {
    private toEntryFormValues(input: AccessListEntrySet): AccessListEntrySetFormValues {
        const { priority, outcome } = input
        const sourceAddresses = input.sourceAddresses.join("\n")

        return {
            priority,
            outcome,
            sourceAddresses,
        }
    }

    private toEntryRequest(input: AccessListEntrySetFormValues): AccessListEntrySet {
        const { priority, outcome } = input
        const sourceAddresses = input.sourceAddresses
            .split("\n")
            .map(item => item.trim())
            .filter(item => item.length > 0)

        return {
            priority,
            outcome,
            sourceAddresses,
        }
    }

    toRequest(input: AccessListFormValues): AccessListRequest {
        const { name, forwardAuthenticationHeader, defaultOutcome, satisfyAll, realm, credentials } = input
        const entries = input.entries.map(entry => this.toEntryRequest(entry))

        return {
            name,
            forwardAuthenticationHeader,
            defaultOutcome,
            satisfyAll,
            realm,
            credentials,
            entries,
        }
    }

    toFormValues(input: AccessListResponse): AccessListFormValues {
        const { name, forwardAuthenticationHeader, defaultOutcome, satisfyAll, realm, credentials } = input
        const entries = input.entries?.map(entry => this.toEntryFormValues(entry)) ?? []

        return {
            name,
            forwardAuthenticationHeader,
            defaultOutcome,
            realm,
            satisfyAll,
            credentials: credentials ?? [],
            entries,
        }
    }
}

export default new AccessListConverter()
