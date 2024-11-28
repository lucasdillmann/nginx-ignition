import CertificateService from "../CertificateService";
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation";
import Notification from "../../../core/components/notification/Notification";
import ReloadNginxAction from "../../nginx/actions/ReloadNginxAction";
import {UnexpectedResponseError} from "../../../core/apiclient/ApiResponse";
import {RenewCertificateResponse} from "../model/RenewCertificateResponse";

class RenewCertificateAction {
    private readonly service: CertificateService

    constructor() {
        this.service = new CertificateService()
    }

    private async invokeCertificateRenew(certificateId: string): Promise<void> {
        return this.service
            .renew(certificateId)
            .then(() => Notification.success(
                `Certificate renewed`,
                `The certificate was renewed successfully`,
            ))
            .then(() => ReloadNginxAction.execute())
            .catch((error: UnexpectedResponseError<RenewCertificateResponse>) => Notification.error(
                `Unable to renew the certificate`,
                error.response.body?.errorReason ??
                `An unexpected error was found while trying to renew the certificate. Please try again later.`,
            ))
    }

    async execute(certificateId: string): Promise<void> {
        return UserConfirmation.askWithCallback(
            `Renewing the certificate can take several seconds and is only recommended when something is wrong with it
            since, by default, nginx ignition will renew it automatically when it's close to expiring. Continue anyway?`,
            () => this.invokeCertificateRenew(certificateId),
        )
    }
}

// eslint-disable-next-line import/no-anonymous-default-export
export default new RenewCertificateAction()
