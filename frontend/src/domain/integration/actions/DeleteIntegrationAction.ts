import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import Notification from "../../../core/components/notification/Notification"
import ReloadNginxAction from "../../nginx/actions/ReloadNginxAction"
import IntegrationService from "../IntegrationService"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { raw } from "../../../core/i18n/I18n"

class DeleteIntegrationAction {
    private readonly service: IntegrationService

    constructor() {
        this.service = new IntegrationService()
    }

    async execute(integrationId: string): Promise<void> {
        return UserConfirmation.ask(MessageKey.FrontendIntegrationDeleteConfirmation)
            .then(() => this.service.delete(integrationId))
            .then(() => {
                Notification.success(
                    { id: MessageKey.CommonTypeDeleted, params: { type: MessageKey.CommonEntityIntegration } },
                    MessageKey.CommonSuccessMessage,
                )
                ReloadNginxAction.execute()
            })
            .catch(error => {
                const message = error?.response?.body?.message
                    ? raw(error.response.body.message)
                    : MessageKey.CommonUnexpectedErrorTryAgain

                Notification.error(
                    { id: MessageKey.CommonUnableToDelete, params: { type: MessageKey.CommonEntityIntegration } },
                    message,
                )
            })
    }
}

export default new DeleteIntegrationAction()
