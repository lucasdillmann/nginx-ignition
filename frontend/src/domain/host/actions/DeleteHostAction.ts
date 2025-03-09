import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import Notification from "../../../core/components/notification/Notification"
import HostService from "../HostService"
import ReloadNginxAction from "../../nginx/actions/ReloadNginxAction"

class DeleteHostAction {
    private readonly service: HostService

    constructor() {
        this.service = new HostService()
    }

    async execute(hostId: string): Promise<void> {
        return UserConfirmation.ask("Do you really want to delete the host?")
            .then(() => this.service.delete(hostId))
            .then(() => Notification.success(`Host deleted`, `The host was deleted successfully`))
            .then(() => ReloadNginxAction.execute())
            .catch(() =>
                Notification.error(
                    `Unable to delete the host`,
                    `An unexpected error was found while trying to delete the host. Please try again later.`,
                ),
            )
    }
}

export default new DeleteHostAction()
