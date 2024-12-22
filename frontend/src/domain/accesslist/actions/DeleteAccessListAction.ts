import AccessListService from "../AccessListService"
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import Notification from "../../../core/components/notification/Notification"
import { UnexpectedResponseError } from "../../../core/apiclient/ApiResponse"

class DeleteAccessListAction {
    private readonly service: AccessListService

    constructor() {
        this.service = new AccessListService()
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const message = error.response?.body?.message
            if (typeof message === "string") {
                Notification.error(`Unable to delete the access list`, message)
                return
            }
        }

        Notification.error(
            `Unable to delete the access list`,
            `An unexpected error was found while trying to delete the access list. Please try again later.`,
        )
    }

    async execute(userId: string): Promise<void> {
        return UserConfirmation.ask("Do you really want to delete the access list?")
            .then(() => this.service.delete(userId))
            .then(() => Notification.success(`Access list deleted`, `The access list was deleted successfully`))
            .catch(error => this.handleError(error))
    }
}

// eslint-disable-next-line import/no-anonymous-default-export
export default new DeleteAccessListAction()
