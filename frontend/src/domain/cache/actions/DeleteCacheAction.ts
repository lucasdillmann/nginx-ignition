import CacheService from "../CacheService"
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import Notification from "../../../core/components/notification/Notification"
import { UnexpectedResponseError } from "../../../core/apiclient/ApiResponse"

class DeleteCacheAction {
    private readonly service: CacheService

    constructor() {
        this.service = new CacheService()
    }

    private handleError(error: Error) {
        if (error instanceof UnexpectedResponseError) {
            const message = error.response?.body?.message
            if (typeof message === "string") {
                Notification.error(`Unable to delete the cache configuration`, message)
                return
            }
        }

        Notification.error(
            `Unable to delete the cache configuration`,
            `An unexpected error was found while trying to delete the cache configuration. Please try again later.`,
        )
    }

    async execute(userId: string): Promise<void> {
        return UserConfirmation.ask("Do you really want to delete the cache configuration?")
            .then(() => this.service.delete(userId))
            .then(() =>
                Notification.success(`Cache configuration deleted`, `The cache configuration was deleted successfully`),
            )
            .catch(error => this.handleError(error))
    }
}

export default new DeleteCacheAction()
