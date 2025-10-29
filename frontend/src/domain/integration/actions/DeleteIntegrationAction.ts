import UserConfirmation from "../../../core/components/confirmation/UserConfirmation"
import Notification from "../../../core/components/notification/Notification"
import ReloadNginxAction from "../../nginx/actions/ReloadNginxAction"
import IntegrationService from "../IntegrationService"

class DeleteIntegrationAction {
    private readonly service: IntegrationService

    constructor() {
        this.service = new IntegrationService()
    }

    async execute(integrationId: string): Promise<void> {
        return UserConfirmation.ask("Do you really want to delete the integration?")
            .then(() => this.service.delete(integrationId))
            .then(() => {
                Notification.success(`Integration deleted`, `The integration was deleted successfully`)
                ReloadNginxAction.execute()
            })
            .catch(error => {
                let message = `An unexpected error was found while trying to delete the integration. Please try again later.`
                if (error && error.response?.body?.message) message = error.response.body.message

                Notification.error(`Unable to delete the integration`, message)
            })
    }
}

export default new DeleteIntegrationAction()
