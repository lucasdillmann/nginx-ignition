import UserService from "../UserService"
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import Notification from "../../../core/components/notification/Notification"
import { UnexpectedResponseError } from "../../../core/apiclient/ApiResponse"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { I18nMessage, raw } from "../../../core/i18n/I18n"

class DeleteAPITokenAction {
    private readonly service: UserService

    constructor() {
        this.service = new UserService()
    }

    private handleError(error: Error) {
        const title = MessageKey.FrontendUserTokensRevokedTitle
        let message: I18nMessage = MessageKey.CommonUnexpectedErrorTryAgain

        if (error instanceof UnexpectedResponseError) {
            const responseMessage = error.response?.body?.message
            if (typeof responseMessage === "string") {
                message = raw(responseMessage)
            }
        }

        Notification.error(title, message)
    }

    async execute(tokenId: string): Promise<void> {
        return UserConfirmation.ask(MessageKey.FrontendUserTokensDeleteConfirmation)
            .then(() => this.service.deleteToken(tokenId))
            .then(() =>
                Notification.success(MessageKey.FrontendUserTokensRevokedTitle, MessageKey.CommonSuccessMessage),
            )
            .catch(error => this.handleError(error))
    }
}

export default new DeleteAPITokenAction()
