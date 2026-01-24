import CertificateService from "../CertificateService"
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import Notification from "../../../core/components/notification/Notification"
import { UnexpectedResponseError } from "../../../core/apiclient/ApiResponse"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { raw } from "../../../core/i18n/I18n"

class DeleteCertificateAction {
    private readonly service: CertificateService

    constructor() {
        this.service = new CertificateService()
    }

    private handleError(error: Error) {
        const title = {
            id: MessageKey.CommonUnableToDelete,
            params: { type: MessageKey.CommonEntityCertificate },
        }
        if (error instanceof UnexpectedResponseError) {
            const message = error.response?.body?.message
            if (typeof message === "string") {
                Notification.error(title, raw(message))
                return
            }
        }

        Notification.error(title, MessageKey.CommonUnexpectedErrorTryAgain)
    }

    async execute(certificateId: string): Promise<void> {
        return UserConfirmation.ask("Do you really want to delete the certificate?")
            .then(() => this.service.delete(certificateId))
            .then(() =>
                Notification.success(
                    { id: MessageKey.CommonTypeDeleted, params: { type: MessageKey.CommonEntityCertificate } },
                    MessageKey.CommonSuccessMessage,
                ),
            )
            .catch(error => this.handleError(error))
    }
}

export default new DeleteCertificateAction()
