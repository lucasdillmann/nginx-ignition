import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import Notification from "../../../core/components/notification/Notification"
import ReloadNginxAction from "../../nginx/actions/ReloadNginxAction"
import VpnService from "../VpnService"

class DeleteVpnAction {
    private readonly service: VpnService

    constructor() {
        this.service = new VpnService()
    }

    async execute(vpnId: string): Promise<void> {
        return UserConfirmation.ask("Do you really want to delete the VPN connection?")
            .then(() => this.service.delete(vpnId))
            .then(() => {
                Notification.success(`VPN connection deleted`, `The VPN connection was deleted successfully`)
                ReloadNginxAction.execute()
            })
            .catch(error => {
                let message = `An unexpected error was found while trying to delete the VPN connection. Please try again later.`
                if (error?.response?.body?.message) message = error.response.body.message

                Notification.error(`Unable to delete the VPN connection`, message)
            })
    }
}

export default new DeleteVpnAction()
