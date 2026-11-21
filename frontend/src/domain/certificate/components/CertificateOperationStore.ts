import {
    CertificateOperationLog,
    CertificateOperationResult,
    CertificateOperationValidation,
} from "../model/CertificateOperationEvent"

export type CertificateOperationKind = "issue" | "renew"

export type CertificateOperationPhase = "running" | "canceling" | "success" | "failure" | "canceled"

export interface CertificateOperationLogLine extends CertificateOperationLog {
    id: number
}

export interface CertificateOperationState {
    kind: CertificateOperationKind
    phase: CertificateOperationPhase
    logs: CertificateOperationLogLine[]
    showDebug: boolean
    truncated: boolean
    logCount: number
    certificateId?: string
    errorReason?: string
}

const maximumLogLines = 5000

export class CertificateOperationStore {
    private readonly listeners: Set<() => void>
    private readonly logs: CertificateOperationLogLine[]
    private cancelHandler?: () => void
    private nextLogId = 0
    private truncated = false
    private state: CertificateOperationState

    constructor(kind: CertificateOperationKind) {
        this.listeners = new Set()
        this.logs = []
        this.state = {
            kind,
            phase: "running",
            logs: this.logs,
            showDebug: false,
            truncated: false,
            logCount: 0,
        }
    }

    getState(): CertificateOperationState {
        return this.state
    }

    subscribe(listener: () => void): () => void {
        this.listeners.add(listener)

        return () => this.listeners.delete(listener)
    }

    setCancelHandler(handler: () => void) {
        this.cancelHandler = handler
    }

    cancel() {
        if (this.state.phase !== "running") return

        this.update({ phase: "canceling" })
        this.cancelHandler?.()
    }

    addLog(log: CertificateOperationLog) {
        this.logs.push({ ...log, id: this.nextLogId++ })

        if (this.logs.length > maximumLogLines) {
            this.logs.splice(0, this.logs.length - maximumLogLines)
            this.truncated = true
        }

        this.update({ logs: this.logs, truncated: this.truncated, logCount: this.state.logCount + 1 })
    }

    setShowDebug(showDebug: boolean) {
        this.update({ showDebug })
    }

    setResult(result: CertificateOperationResult) {
        if (result.success) {
            this.update({
                phase: "success",
                certificateId: result.certificateId,
                errorReason: undefined,
            })
            return
        }

        this.update({ phase: "failure", errorReason: result.errorReason })
    }

    setValidation(validation: CertificateOperationValidation) {
        this.update({ phase: "failure", errorReason: validation.message })
    }

    setFailure(errorReason: string) {
        this.update({ phase: "failure", errorReason })
    }

    setCanceled() {
        this.update({ phase: "canceled" })
    }

    private update(patch: Partial<CertificateOperationState>) {
        this.state = { ...this.state, ...patch }
        this.listeners.forEach(listener => listener())
    }
}
