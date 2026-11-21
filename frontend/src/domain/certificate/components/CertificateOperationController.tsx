import React from "react"
import { themedModal } from "../../../core/components/theme/ThemedResources"
import { I18n } from "../../../core/i18n/I18n"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import CertificateService, { CertificateOperationListeners } from "../CertificateService"
import { UnexpectedResponseError } from "../../../core/apiclient/ApiResponse"
import { IssueCertificateRequest } from "../model/IssueCertificateRequest"
import { CertificateOperationResult, CertificateOperationValidation } from "../model/CertificateOperationEvent"
import { CertificateOperationKind, CertificateOperationStore } from "./CertificateOperationStore"
import CertificateOperationModal from "./CertificateOperationModal"

export interface CertificateOperationOptions {
    onValidation?: (validation: CertificateOperationValidation) => void
}

type OperationRunner = (
    listeners: CertificateOperationListeners,
    signal: AbortSignal,
) => Promise<CertificateOperationResult>

class CertificateOperationController {
    private readonly service: CertificateService

    constructor() {
        this.service = new CertificateService()
    }

    async issue(
        request: IssueCertificateRequest,
        options: CertificateOperationOptions = {},
    ): Promise<CertificateOperationResult> {
        return this.run("issue", options, (listeners, signal) => this.service.issue(request, listeners, signal))
    }

    async renew(id: string, options: CertificateOperationOptions = {}): Promise<CertificateOperationResult> {
        return this.run("renew", options, (listeners, signal) => this.service.renew(id, listeners, signal))
    }

    private async run(
        kind: CertificateOperationKind,
        options: CertificateOperationOptions,
        execute: OperationRunner,
    ): Promise<CertificateOperationResult> {
        const store = new CertificateOperationStore(kind)
        const abortController = new AbortController()

        const title =
            kind === "issue"
                ? MessageKey.FrontendCertificateOperationIssueTitle
                : MessageKey.FrontendCertificateOperationRenewTitle
        const modal: { instance?: { destroy: () => void } } = {}
        const close = () => modal.instance?.destroy()

        modal.instance = themedModal().info({
            title: <I18n id={title} />,
            content: <CertificateOperationModal store={store} onClose={close} />,
            footer: null,
            icon: null,
            closable: false,
            mask: { closable: false },
            keyboard: false,
            width: 780,
        })

        store.setCancelHandler(() => abortController.abort())

        try {
            return await execute(
                {
                    onLog: log => store.addLog(log),
                    onValidation: validation => {
                        if (kind === "issue") {
                            options.onValidation?.(validation)
                            close()
                            return
                        }

                        validation.consistencyProblems?.forEach(problem => {
                            store.addLog({
                                time: new Date().toISOString(),
                                level: "error",
                                message: `${problem.fieldPath}: ${problem.message}`,
                            })
                        })
                        store.setValidation(validation)
                    },
                    onResult: result => store.setResult(result),
                    onError: message => store.setFailure(message),
                },
                abortController.signal,
            )
        } catch (error) {
            if (isAbortError(error)) {
                store.setCanceled()

                return { success: false }
            }

            const message = resolveErrorMessage(error)
            store.setFailure(message)

            return { success: false, errorReason: message }
        }
    }
}

function isAbortError(error: unknown): boolean {
    return (
        (error instanceof DOMException && error.name === "AbortError") ||
        (error instanceof Error && error.name === "AbortError")
    )
}

function resolveErrorMessage(error: unknown): string {
    if (error instanceof UnexpectedResponseError) {
        const message = error.response.body?.message
        if (typeof message === "string" && message.length > 0) return message

        return `Request failed with status ${error.response.statusCode}`
    }

    if (error instanceof Error && error.message.length > 0) return error.message

    return String(error)
}

export default new CertificateOperationController()
