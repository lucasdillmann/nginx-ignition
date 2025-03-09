import CertificateService from "../CertificateService"
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import Notification from "../../../core/components/notification/Notification"
import { UnexpectedResponseError } from "../../../core/apiclient/ApiResponse"

class DeleteCertificateAction {
    private readonly service: CertificateService

    constructor() {
        this.service = new CertificateService()
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const message = error.response?.body?.message
            if (typeof message === "string") {
                Notification.error(`Unable to delete the certificate`, message)
                return
            }
        }

        Notification.error(
            `Unable to delete the certificate`,
            `An unexpected error was found while trying to delete the certificate. Please try again later.`,
        )
    }

    async execute(certificateId: string): Promise<void> {
        return UserConfirmation.ask("Do you really want to delete the certificate?")
            .then(() => this.service.delete(certificateId))
            .then(() => Notification.success(`Certificate deleted`, `The certificate was deleted successfully`))
            .catch(error => this.handleError(error))
    }
}

export default new DeleteCertificateAction()
