import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import Notification from "../../../core/components/notification/Notification"
import HostService from "../HostService"
import ReloadNginxAction from "../../nginx/actions/ReloadNginxAction"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { i18n } from "../../../core/i18n/I18n"

class DeleteHostAction {
    private readonly service: HostService

    constructor() {
        this.service = new HostService()
    }

    async execute(hostId: string): Promise<void> {
        return UserConfirmation.ask("Do you really want to delete the host?")
            .then(() => this.service.delete(hostId))
            .then(() => {
                Notification.success(
                    { id: MessageKey.CommonTypeDeleted, params: { type: i18n(MessageKey.CommonEntityHost) } },
                    MessageKey.CommonSuccessMessage,
                )
                ReloadNginxAction.execute()
            })
            .catch(() =>
                Notification.error(
                    { id: MessageKey.CommonUnableToDelete, params: { type: i18n(MessageKey.CommonEntityHost) } },
                    MessageKey.CommonUnexpectedErrorTryAgain,
                ),
            )
    }
}

export default new DeleteHostAction()
