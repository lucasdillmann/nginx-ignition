import React from "react"
import { App, ConfigProvider } from "antd"
import ErrorBoundary from "../core/components/errorboundary/ErrorBoundary"
import AppContainer from "./AppContainer"
import ApiClientEventDispatcher from "../core/apiclient/event/ApiClientEventDispatcher"
import AuthenticationTokenApiClientEventListener from "../core/authentication/AuthenticationTokenApiClientEventListener"
import SessionExpiredApiClientEventListener from "../core/authentication/SessionExpiredApiClientEventListener"
import ThemeContext from "../core/components/context/ThemeContext"
import ThemedResources from "../core/components/theme/ThemedResources"

interface NginxIgnitionState {
    error?: Error
    darkMode: boolean
}

export default class NginxIgnition extends React.Component<unknown, NginxIgnitionState> {
    constructor(props: unknown) {
        super(props)

        this.state = {
            darkMode: ThemeContext.isDarkMode(),
        }
    }

    private handleThemeChange(darkMode: boolean) {
        this.setState({ darkMode })
    }

    componentDidMount() {
        ThemeContext.register(this.handleThemeChange.bind(this))
        ApiClientEventDispatcher.register(new AuthenticationTokenApiClientEventListener())
        ApiClientEventDispatcher.register(new SessionExpiredApiClientEventListener())

        const preloader = document.getElementById("preloader") as HTMLElement
        preloader?.remove()
    }

    componentWillUnmount() {
        ThemeContext.deregister(this.handleThemeChange.bind(this))
    }

    render() {
        document.documentElement.setAttribute("data-theme", ThemeContext.isDarkMode() ? "dark" : "light")

        return (
            <ErrorBoundary>
                <ConfigProvider theme={{ algorithm: ThemeContext.algorithm() }}>
                    <App>
                        <ThemedResources />
                        <AppContainer />
                    </App>
                </ConfigProvider>
            </ErrorBoundary>
        )
    }
}
