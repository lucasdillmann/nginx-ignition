export interface ServerCertificateResponse {
    id: string
    domainNames: string[]
    providerId: string
    issuedAt: string
    validUntil: string
    validFrom: string
    renewAfter?: string
    parameters: Record<string, any>
}
