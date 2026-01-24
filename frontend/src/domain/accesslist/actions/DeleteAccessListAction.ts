import AccessListService from "../AccessListService"
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import Notification from "../../../core/components/notification/Notification"
import { UnexpectedResponseError } from "../../../core/apiclient/ApiResponse"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { raw } from "../../../core/i18n/I18n"

class DeleteAccessListAction {
    private readonly service: AccessListService

    constructor() {
        this.service = new AccessListService()
    }

    private handleError(error: Error) {
        const title = {
            id: MessageKey.CommonUnableToDelete,
            params: {
                type: MessageKey.CommonEntityAccessList,
            },
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

    async execute(userId: string): Promise<void> {
        return UserConfirmation.ask("Do you really want to delete the access list?")
            .then(() => this.service.delete(userId))
            .then(() =>
                Notification.success(
                    { id: MessageKey.CommonTypeDeleted, params: { type: MessageKey.CommonEntityAccessList } },
                    MessageKey.CommonSuccessMessage,
                ),
            )
            .catch(error => this.handleError(error))
    }
}

export default new DeleteAccessListAction()
