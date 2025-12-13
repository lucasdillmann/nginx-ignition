import ServerCertificateService from "../ServerCertificateService"
import UserConfirmation from "../../../../core/components/confirmation/UserConfirmation"
import Notification from "../../../../core/components/notification/Notification"
import ReloadNginxAction from "../../../nginx/actions/ReloadNginxAction"
import { UnexpectedResponseError } from "../../../../core/apiclient/ApiResponse"
import { RenewServerCertificateResponse } from "../model/RenewServerCertificateResponse"

class RenewServerCertificateAction {
    private readonly service: ServerCertificateService

    constructor() {
        this.service = new ServerCertificateService()
    }

    private async invokeCertificateRenew(certificateId: string): Promise<void> {
        return this.service
            .renew(certificateId)
            .then(() => {
                Notification.success(`Server certificate renewed`, `The server certificate was renewed successfully`)
                ReloadNginxAction.execute()
            })
            .catch((error: UnexpectedResponseError<RenewServerCertificateResponse>) =>
                Notification.error(
                    `Unable to renew the server certificate`,
                    error.response.body?.errorReason ??
                        `An unexpected error was found while trying to renew the server certificate. Please try again later.`,
                ),
            )
    }

    async execute(certificateId: string): Promise<void> {
        return UserConfirmation.askWithCallback(
            `Renewing the server certificate can take several seconds, and is recommended only when something is wrong with it
             (nginx ignition will renew it automatically when it's close to expiring). Continue anyway?`,
            () => this.invokeCertificateRenew(certificateId),
        )
    }
}

export default new RenewServerCertificateAction()
