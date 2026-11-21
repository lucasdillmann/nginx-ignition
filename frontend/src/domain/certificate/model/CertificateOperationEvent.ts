import { SseEvent } from "../../../core/apiclient/SseReader"

export type CertificateOperationLogLevel = "debug" | "info" | "warn" | "error"

export interface CertificateOperationLog {
    time: string
    level: CertificateOperationLogLevel
    message: string
}

export interface CertificateOperationResult {
    success: boolean
    errorReason?: string
    certificateId?: string
}

export interface CertificateOperationConsistencyProblem {
    fieldPath: string
    message: string
}

export interface CertificateOperationValidation {
    message: string
    consistencyProblems?: CertificateOperationConsistencyProblem[]
}

export type CertificateOperationEvent =
    | { type: "log"; log: CertificateOperationLog }
    | { type: "result"; result: CertificateOperationResult }
    | { type: "validation"; validation: CertificateOperationValidation }
    | { type: "error"; message: string }

export function parseCertificateOperationEvent(event: SseEvent): CertificateOperationEvent | null {
    let payload: any

    try {
        payload = JSON.parse(event.data)
    } catch (_) {
        return null
    }

    switch (event.event) {
        case "log":
            return { type: "log", log: payload }
        case "result":
            return { type: "result", result: payload }
        case "validation":
            return { type: "validation", validation: payload }
        case "error":
            return { type: "error", message: payload?.message ?? String(payload) }
        default:
            return null
    }
}
