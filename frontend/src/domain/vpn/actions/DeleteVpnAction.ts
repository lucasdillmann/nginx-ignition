import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import Notification from "../../../core/components/notification/Notification"
import ReloadNginxAction from "../../nginx/actions/ReloadNginxAction"
import VpnService from "../VpnService"
import MessageKey from "../../../core/i18n/model/MessageKey.generated"
import { raw } from "../../../core/i18n/I18n"

class DeleteVpnAction {
    private readonly service: VpnService

    constructor() {
        this.service = new VpnService()
    }

    async execute(vpnId: string): Promise<void> {
        return UserConfirmation.ask(MessageKey.FrontendVpnDeleteConfirmation)
            .then(() => this.service.delete(vpnId))
            .then(() => {
                Notification.success(
                    { id: MessageKey.CommonTypeDeleted, params: { type: MessageKey.CommonVpnConnection } },
                    MessageKey.CommonSuccessMessage,
                )
                ReloadNginxAction.execute()
            })
            .catch(error => {
                const message = error?.response?.body?.message
                    ? raw(error.response.body.message)
                    : MessageKey.CommonUnexpectedErrorTryAgain

                Notification.error(
                    {
                        id: MessageKey.CommonUnableToDelete,
                        params: { type: MessageKey.CommonVpnConnection },
                    },
                    message,
                )
            })
    }
}

export default new DeleteVpnAction()
