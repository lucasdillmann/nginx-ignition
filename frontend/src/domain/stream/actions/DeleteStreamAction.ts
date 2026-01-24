import StreamService from "../StreamService"
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import Notification from "../../../core/components/notification/Notification"
import { UnexpectedResponseError } from "../../../core/apiclient/ApiResponse"
import ReloadNginxAction from "../../nginx/actions/ReloadNginxAction"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { i18n, raw } from "../../../core/i18n/I18n"

class DeleteStreamAction {
    private readonly service: StreamService

    constructor() {
        this.service = new StreamService()
    }

    private handleError(error: Error) {
        const title = { id: MessageKey.CommonUnableToDelete, params: { type: i18n(MessageKey.CommonEntityStream) } }
        if (error instanceof UnexpectedResponseError) {
            const message = error.response?.body?.message
            if (typeof message === "string") {
                Notification.error(title, raw(message))
                return
            }
        }

        Notification.error(title, MessageKey.CommonUnexpectedErrorTryAgain)
    }

    async execute(streamId: string): Promise<void> {
        return UserConfirmation.ask("Do you really want to delete the stream?")
            .then(() => this.service.delete(streamId))
            .then(() =>
                Notification.success(
                    { id: MessageKey.CommonTypeDeleted, params: { type: i18n(MessageKey.CommonEntityStream) } },
                    MessageKey.CommonSuccessMessage,
                ),
            )
            .then(() => ReloadNginxAction.execute())
            .catch(error => this.handleError(error))
    }
}

export default new DeleteStreamAction()
