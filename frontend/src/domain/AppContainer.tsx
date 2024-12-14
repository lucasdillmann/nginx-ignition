import React from "react"
import AppContext, { loadAppContextData } from "../core/components/context/AppContext"
import AppRouter from "../core/components/router/AppRouter"
import Routes from "./Routes"
import FullPagePreloader from "../core/components/preloader/FullPagePreloader"
import FullPageError from "../core/components/error/FullPageError"
import ShellUserMenu from "./user/components/ShellUserMenu"
import NginxControl from "./nginx/components/NginxControl"
import CommonNotifications from "../core/components/notification/CommonNotifications"

interface AppContainerState {
    loading: boolean
    error?: Error
}

export default class AppContainer extends React.Component<unknown, AppContainerState> {
    constructor(props: any) {
        super(props)
        this.state = {
            loading: true,
        }
    }

    componentDidMount() {
        loadAppContextData()
            .then(context => {
                AppContext.replace(context)
                this.setState({ loading: false })
            })
            .catch(error => {
                CommonNotifications.failedToFetch()
                this.setState({ error, loading: false })
            })
    }

    render() {
        const { error, loading } = this.state
        if (error !== undefined) return <FullPageError error={error} />
        if (loading) return <FullPagePreloader />

        return <AppRouter routes={Routes} userMenu={<ShellUserMenu />} serverControl={<NginxControl />} />
    }
}
