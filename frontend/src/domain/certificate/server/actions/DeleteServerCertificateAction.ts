import ServerCertificateService from "../ServerCertificateService"
import UserConfirmation from "../../../../core/components/confirmation/UserConfirmation"
import Notification from "../../../../core/components/notification/Notification"
import { UnexpectedResponseError } from "../../../../core/apiclient/ApiResponse"

class DeleteServerCertificateAction {
    private readonly service: ServerCertificateService

    constructor() {
        this.service = new ServerCertificateService()
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const message = error.response?.body?.message
            if (typeof message === "string") {
                Notification.error(`Unable to delete the server certificate`, message)
                return
            }
        }

        Notification.error(
            `Unable to delete the server certificate`,
            `An unexpected error was found while trying to delete the server certificate. Please try again later.`,
        )
    }

    async execute(serverCertificateId: string): Promise<void> {
        return UserConfirmation.ask("Do you really want to delete the server certificate?")
            .then(() => this.service.delete(serverCertificateId))
            .then(() =>
                Notification.success(`Server certificate deleted`, `The server certificate was deleted successfully`),
            )
            .catch(error => this.handleError(error))
    }
}

export default new DeleteServerCertificateAction()
