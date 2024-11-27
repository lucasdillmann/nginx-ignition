import UserConfirmation from "../../../core/components/confirmation/UserConfirmation";
import Notification from "../../../core/components/notification/Notification";
import HostService from "../HostService";
import NginxReload from "../../../core/components/nginx/NginxReload";

class DeleteHostAction {
    private readonly service: HostService

    constructor() {
        this.service = new HostService()
    }

    async execute(hostId: string): Promise<void> {
        return UserConfirmation
            .ask("Do you really want to delete the host?")
            .then(() => this.service.delete(hostId))
            .then(() => Notification.success(
                `Host deleted`,
                `The host was deleted successfully`,
            ))
            .then(() => NginxReload.ask())
            .catch(() => Notification.error(
                `Unable to delete the host`,
                `An unexpected error was found while trying to delete the host. Please try again later.`,
            ))
    }
}

// eslint-disable-next-line import/no-anonymous-default-export
export default new DeleteHostAction()
