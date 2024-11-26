import CertificateService from "../CertificateService";
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation";
import Notification from "../../../core/components/notification/Notification";
import NginxReload from "../../../core/components/nginx/NginxReload";

class DeleteCertificateAction {
    private readonly service: CertificateService

    constructor() {
        this.service = new CertificateService()
    }

    async execute(certificateId: string): Promise<void> {
        return UserConfirmation
            .ask("Do you really want to delete the certificate?")
            .then(() => this.service.delete(certificateId))
            .then(() => Notification.success(
                `Certificate deleted`,
                `The certificate was deleted successfully`,
            ))
            .then(() => NginxReload.ask())
            .catch(() => Notification.error(
                `Unable to delete the certificate`,
                `An unexpected error was found while trying to delete the certificate. Please try again later.`,
            ))
    }
}

// eslint-disable-next-line import/no-anonymous-default-export
export default new DeleteCertificateAction()
