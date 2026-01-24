import CertificateService from "../CertificateService"
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import Notification from "../../../core/components/notification/Notification"
import ReloadNginxAction from "../../nginx/actions/ReloadNginxAction"
import { UnexpectedResponseError } from "../../../core/apiclient/ApiResponse"
import { RenewCertificateResponse } from "../model/RenewCertificateResponse"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { raw } from "../../../core/i18n/I18n"

class RenewCertificateAction {
    private readonly service: CertificateService

    constructor() {
        this.service = new CertificateService()
    }

    private async invokeCertificateRenew(certificateId: string): Promise<void> {
        return this.service
            .renew(certificateId)
            .then(() => {
                Notification.success(MessageKey.FrontendCertificateRenewed, MessageKey.CommonSuccessMessage)
                ReloadNginxAction.execute()
            })
            .catch((error: UnexpectedResponseError<RenewCertificateResponse>) =>
                Notification.error(
                    MessageKey.FrontendCertificateRenewFailed,
                    error.response.body?.errorReason
                        ? raw(error.response.body?.errorReason)
                        : MessageKey.CommonUnexpectedErrorTryAgain,
                ),
            )
    }

    async execute(certificateId: string): Promise<void> {
        return UserConfirmation.askWithCallback(
            `Renewing the certificate can take several seconds and is only recommended when something is wrong with it
            since, by default, nginx ignition will renew it automatically when it's close to expiring. Continue anyway?`,
            () => this.invokeCertificateRenew(certificateId),
        )
    }
}

export default new RenewCertificateAction()
