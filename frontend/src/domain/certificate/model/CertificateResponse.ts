export interface CertificateResponse {
    id: string
    domainNames: string[]
    providerId: string
    issuedAt: string
    validUntil: string
    validFrom: string
    renewAfter?: string
}
