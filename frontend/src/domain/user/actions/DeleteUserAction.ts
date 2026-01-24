import UserService from "../UserService"
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import Notification from "../../../core/components/notification/Notification"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { i18n } from "../../../core/i18n/I18n"

class DeleteUserAction {
    private readonly service: UserService

    constructor() {
        this.service = new UserService()
    }

    async execute(userId: string): Promise<void> {
        return UserConfirmation.ask("Do you really want to delete the user?")
            .then(() => this.service.delete(userId))
            .then(() =>
                Notification.success(
                    { id: MessageKey.CommonTypeDeleted, params: { type: i18n(MessageKey.CommonEntityUser) } },
                    MessageKey.CommonSuccessMessage,
                ),
            )
            .catch(() =>
                Notification.error(
                    { id: MessageKey.CommonUnableToDelete, params: { type: i18n(MessageKey.CommonEntityUser) } },
                    MessageKey.CommonUnexpectedErrorTryAgain,
                ),
            )
    }
}

export default new DeleteUserAction()
