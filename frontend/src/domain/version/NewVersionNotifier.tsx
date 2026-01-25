import AppContext from "../../core/components/context/AppContext"
import LocalStorageRepository from "../../core/repository/LocalStorageRepository"
import { themedModal } from "../../core/components/theme/ThemedResources"
import { i18n } from "../../core/i18n/I18n"
import MessageKey from "../../core/i18n/model/MessageKey.generated"

class NewVersionNotifier {
    private readonly repository: LocalStorageRepository<string>

    constructor() {
        this.repository = new LocalStorageRepository("nginxIgnition.lastNotifiedNewVersion")
    }

    checkAndNotify() {
        const { version } = AppContext.get().configuration
        const { current, latest } = version
        const lastNotifiedNewVersion = this.repository.get()

        if (current && latest && current !== latest && lastNotifiedNewVersion !== latest) {
            this.repository.set(latest)
            const modalInstance = themedModal().info({
                type: "info",
                title: i18n(MessageKey.FrontendVersionNotifierTitle),
                content: i18n({ id: MessageKey.FrontendVersionNotifierContent, params: { version: latest } }),
                okText: i18n(MessageKey.FrontendVersionNotifierOpenRelease),
                closable: true,
                width: 600,
                okButtonProps: {
                    onClick: () => {
                        window.open(`https://github.com/lucasdillmann/nginx-ignition/releases/${latest}`, "_blank")
                        modalInstance.destroy()
                    },
                },
            })
        }
    }
}

export default new NewVersionNotifier()
