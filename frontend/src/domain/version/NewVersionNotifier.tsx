import AppContext from "../../core/components/context/AppContext"
import LocalStorageRepository from "../../core/repository/LocalStorageRepository"
import { themedModal } from "../../core/components/theme/ThemedResources"

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
                title: "New version available",
                content: `nginx ignition ${latest} is available with (probably) new features. You can check the project release page for the changelog and more.`,
                okText: "Open release page",
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
