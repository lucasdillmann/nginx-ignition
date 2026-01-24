import CacheService from "../CacheService"
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import Notification from "../../../core/components/notification/Notification"
import { UnexpectedResponseError } from "../../../core/apiclient/ApiResponse"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { i18n, I18nMessage, raw } from "../../../core/i18n/I18n"

class DeleteCacheAction {
    private readonly service: CacheService

    constructor() {
        this.service = new CacheService()
    }

    private handleError(error: Error) {
        const title = {
            id: MessageKey.CommonUnableToDelete,
            params: { type: i18n(MessageKey.CommonEntityCacheConfiguration) },
        }
        let message: I18nMessage = MessageKey.CommonUnexpectedErrorTryAgain

        if (error instanceof UnexpectedResponseError) {
            const responseMessage = error.response?.body?.message
            if (typeof responseMessage === "string") {
                message = raw(responseMessage)
            }
        }

        Notification.error(title, message)
    }

    async execute(userId: string): Promise<void> {
        return UserConfirmation.ask("Do you really want to delete the cache configuration?")
            .then(() => this.service.delete(userId))
            .then(() =>
                Notification.success(
                    {
                        id: MessageKey.CommonTypeDeleted,
                        params: { type: i18n(MessageKey.CommonEntityCacheConfiguration) },
                    },
                    MessageKey.CommonSuccessMessage,
                ),
            )
            .catch(error => this.handleError(error))
    }
}

export default new DeleteCacheAction()
