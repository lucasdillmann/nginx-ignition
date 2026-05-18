import React from "react"
import { App, ConfigProvider } from "antd"
import ErrorBoundary from "../core/components/errorboundary/ErrorBoundary"
import AppContainer from "./AppContainer"
import ApiClientEventDispatcher from "../core/apiclient/event/ApiClientEventDispatcher"
import AuthenticationTokenApiClientEventListener from "../core/authentication/AuthenticationTokenApiClientEventListener"
import SessionExpiredApiClientEventListener from "../core/authentication/SessionExpiredApiClientEventListener"
import ThemeContext from "../core/components/context/ThemeContext"
import ThemedResources from "../core/components/theme/ThemedResources"
import { Locale } from "antd/es/locale"
import I18nContext from "../core/i18n/I18nContext"
import { buildI18nLocale } from "../core/i18n/I18nLocale"

interface NginxIgnitionState {
    error?: Error
    darkMode: boolean
    locale: Locale
}

export default class NginxIgnition extends React.Component<unknown, NginxIgnitionState> {
    constructor(props: unknown) {
        super(props)

        this.state = {
            darkMode: ThemeContext.isDarkMode(),
            locale: buildI18nLocale(),
        }
    }

    private handleThemeChange(darkMode: boolean) {
        document.documentElement.setAttribute("data-theme", darkMode ? "dark" : "light")
        this.setState({ darkMode })
    }

    private handleLanguageChange() {
        this.setState({
            locale: buildI18nLocale(),
        })
    }

    componentDidMount() {
        ThemeContext.register(this.handleThemeChange.bind(this))
        I18nContext.register(this.handleLanguageChange.bind(this))
        ApiClientEventDispatcher.register(new AuthenticationTokenApiClientEventListener())
        ApiClientEventDispatcher.register(new SessionExpiredApiClientEventListener())

        document.documentElement.setAttribute("data-theme", ThemeContext.isDarkMode() ? "dark" : "light")
        const preloader = document.getElementById("preloader") as HTMLElement
        preloader?.remove()
    }

    componentWillUnmount() {
        ThemeContext.deregister(this.handleThemeChange.bind(this))
        I18nContext.deregister(this.handleLanguageChange.bind(this))
    }

    render() {
        const { darkMode, locale } = this.state
        const selectTheme = darkMode
            ? {
                  hoverBorderColor: "#6e7d8a",
                  activeBorderColor: "#8797a6",
              }
            : {
                  hoverBorderColor: "#8d9baa",
                  activeBorderColor: "#738293",
              }

        return (
            <ErrorBoundary>
                <ConfigProvider
                    theme={{
                        algorithm: ThemeContext.algorithm(),
                        cssVar: {},
                        components: {
                            Select: selectTheme,
                        },
                    }}
                    locale={locale}
                >
                    <App>
                        <ThemedResources />
                        <AppContainer />
                    </App>
                </ConfigProvider>
            </ErrorBoundary>
        )
    }
}
