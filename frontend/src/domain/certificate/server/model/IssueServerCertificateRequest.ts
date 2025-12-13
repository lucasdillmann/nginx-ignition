export interface IssueServerCertificateRequest {
    providerId: string
    domainNames: string[]
    parameters: Record<string, any>
}
