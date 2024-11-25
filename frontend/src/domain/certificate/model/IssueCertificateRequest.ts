export interface IssueCertificateRequest {
    providerId: string
    domainNames: string[]
    parameters: Record<string, any>
}
