import React from "react"
import AppContext, { loadAppContextData } from "../core/components/context/AppContext"
import AppRouter from "../core/components/router/AppRouter"
import Routes from "./Routes"
import FullPagePreloader from "../core/components/preloader/FullPagePreloader"
import FullPageError from "../core/components/error/FullPageError"
import ShellUserMenu from "./user/components/ShellUserMenu"
import NginxControl from "./nginx/components/NginxControl"
import CommonNotifications from "../core/components/notification/CommonNotifications"
import NewVersionNotifier from "./version/NewVersionNotifier"
import I18nService from "../core/i18n/I18nService"

interface AppContainerState {
    loading: boolean
    error?: Error
}

export default class AppContainer extends React.Component<unknown, AppContainerState> {
    private nativePreloaderPresent = false

    constructor(props: any) {
        super(props)
        this.state = {
            loading: true,
        }
    }

    private async boot() {
        return new I18nService()
            .initContext()
            .then(() => loadAppContextData())
            .then(context => {
                AppContext.replace({
                    ...context,
                    container: this,
                })
                this.setState({ loading: false })
                NewVersionNotifier.checkAndNotify()
            })
            .catch(error => {
                CommonNotifications.failedToFetch()
                this.setState({ error, loading: false })
            })
    }

    async reload() {
        this.setState({ loading: true }, () => this.boot())
    }

    componentDidMount() {
        this.nativePreloaderPresent = document.getElementById("preloader") !== null
        this.reload()
    }

    componentDidUpdate(_prevProps: unknown, prevState: AppContainerState) {
        if (prevState.loading && !this.state.loading && this.nativePreloaderPresent) {
            document.getElementById("preloader")?.remove()
            this.nativePreloaderPresent = false
        }
    }

    render() {
        const { error, loading } = this.state
        if (error !== undefined) return <FullPageError error={error} />
        if (loading) {
            if (this.nativePreloaderPresent) return null
            return <FullPagePreloader />
        }

        return <AppRouter routes={Routes} userMenu={<ShellUserMenu />} serverControl={<NginxControl />} />
    }
}
