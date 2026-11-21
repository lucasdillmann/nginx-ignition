import CertificateGateway from "./CertificateGateway"
import PageResponse from "../../core/pagination/PageResponse"
import { requireNullablePayload, requireSuccessPayload, requireSuccessResponse } from "../../core/apiclient/ApiResponse"
import { CertificateResponse } from "./model/CertificateResponse"
import AvailableProviderResponse from "./model/AvailableProviderResponse"
import { IssueCertificateRequest } from "./model/IssueCertificateRequest"
import {
    CertificateOperationLog,
    CertificateOperationResult,
    CertificateOperationValidation,
    parseCertificateOperationEvent,
} from "./model/CertificateOperationEvent"
import { SseEvent } from "../../core/apiclient/SseReader"

export interface CertificateOperationListeners {
    onLog?: (log: CertificateOperationLog) => void
    onValidation?: (validation: CertificateOperationValidation) => void
    onResult?: (result: CertificateOperationResult) => void
    onError?: (message: string) => void
}

export default class CertificateService {
    private readonly gateway: CertificateGateway

    constructor() {
        this.gateway = new CertificateGateway()
    }

    async list(
        pageSize?: number,
        pageNumber?: number,
        searchTerms?: string,
    ): Promise<PageResponse<CertificateResponse>> {
        return this.gateway.getPage(pageSize, pageNumber, searchTerms).then(requireSuccessPayload)
    }

    async listAll(searchTerms?: string): Promise<CertificateResponse[]> {
        const pageSize = 100
        const certificates: CertificateResponse[] = []
        let pageNumber = 0

        while (true) {
            const page = await this.list(pageSize, pageNumber, searchTerms)
            certificates.push(...page.contents)

            if (page.contents.length === 0 || certificates.length >= page.totalItems) break

            pageNumber++
        }

        return certificates
    }

    async delete(id: string): Promise<void> {
        return this.gateway.delete(id).then(requireSuccessResponse)
    }

    async issue(
        certificate: IssueCertificateRequest,
        listeners: CertificateOperationListeners,
        signal?: AbortSignal,
    ): Promise<CertificateOperationResult> {
        return this.stream(onEvent => this.gateway.issueStream(certificate, onEvent, signal), listeners)
    }

    async renew(
        id: string,
        listeners: CertificateOperationListeners,
        signal?: AbortSignal,
    ): Promise<CertificateOperationResult> {
        return this.stream(onEvent => this.gateway.renewStream(id, onEvent, signal), listeners)
    }

    async availableProviders(): Promise<AvailableProviderResponse[]> {
        return this.gateway.getAvailableProviders().then(requireSuccessPayload)
    }

    async getById(id: string): Promise<CertificateResponse | undefined> {
        return this.gateway.getById(id).then(requireNullablePayload)
    }

    private async stream(
        run: (onEvent: (event: SseEvent) => void) => Promise<void>,
        listeners: CertificateOperationListeners,
    ): Promise<CertificateOperationResult> {
        let result: CertificateOperationResult | undefined

        await run(event => {
            const parsed = parseCertificateOperationEvent(event)
            if (parsed === null) return

            switch (parsed.type) {
                case "log":
                    listeners.onLog?.(parsed.log)
                    break
                case "result":
                    result = parsed.result
                    listeners.onResult?.(parsed.result)
                    break
                case "validation":
                    result = { success: false, errorReason: parsed.validation.message }
                    listeners.onValidation?.(parsed.validation)
                    break
                case "error":
                    result = { success: false, errorReason: parsed.message }
                    listeners.onError?.(parsed.message)
                    break
            }
        })

        if (result === undefined) throw new Error("The certificate operation stream ended without a result")

        return result
    }
}
