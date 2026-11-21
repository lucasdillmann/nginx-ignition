import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import CertificateOperationController from "../components/CertificateOperationController"

class RenewCertificateAction {
    async execute(certificateId: string): Promise<void> {
        return new Promise(resolve => {
            void UserConfirmation.askWithCallback(MessageKey.FrontendCertificateRenewConfirmation, () => {
                void CertificateOperationController.renew(certificateId).finally(resolve)

                return Promise.resolve()
            })
        })
    }
}

export default new RenewCertificateAction()
