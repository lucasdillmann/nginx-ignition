import UserService from "../UserService";
import UserConfirmation from "../../../core/components/confirmation/UserConfirmation";
import Notification from "../../../core/components/notification/Notification";

class DeleteUserAction {
    private readonly service: UserService

    constructor() {
        this.service = new UserService()
    }

    async execute(userId: string): Promise<void> {
        return UserConfirmation
            .ask("Do you really want to delete the user?")
            .then(() => this.service.delete(userId))
            .then(() => Notification.success(
                `User deleted`,
                `The user was deleted successfully`,
            ))
            .catch(() => Notification.error(
                `Unable to delete the user`,
                `An unexpected error was found while trying to delete the user. Please try again later.`,
            ))
    }
}

// eslint-disable-next-line import/no-anonymous-default-export
export default new DeleteUserAction()
