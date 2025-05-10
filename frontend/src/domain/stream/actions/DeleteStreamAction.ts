import StreamService from "../StreamService"
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import Notification from "../../../core/components/notification/Notification"
import { UnexpectedResponseError } from "../../../core/apiclient/ApiResponse"

class DeleteStreamAction {
    private readonly service: StreamService

    constructor() {
        this.service = new StreamService()
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const message = error.response?.body?.message
            if (typeof message === "string") {
                Notification.error(`Unable to delete the stream`, message)
                return
            }
        }

        Notification.error(
            `Unable to delete the stream`,
            `An unexpected error was found while trying to delete the stream. Please try again later.`,
        )
    }

    async execute(streamId: string): Promise<void> {
        return UserConfirmation.ask("Do you really want to delete the stream?")
            .then(() => this.service.delete(streamId))
            .then(() => Notification.success(`Stream deleted`, `The stream was deleted successfully`))
            .catch(error => this.handleError(error))
    }
}

export default new DeleteStreamAction()
