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
import MessageKey from "../core/i18n/model/MessageKey.generated"
import { i18n } from "../core/i18n/I18n"

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
            locale: this.buildLocale(),
        }
    }

    private buildLocale(): Locale {
        const { currentLanguage, defaultLanguage } = I18nContext.get()
        return {
            locale: currentLanguage ?? defaultLanguage,
            Form: {
                optional: `(${i18n(MessageKey.CommonOptional)})`,
                defaultValidateMessages: {},
            },
            Empty: {
                description: i18n(MessageKey.CommonNoData),
            },
        }
    }

    private handleThemeChange(darkMode: boolean) {
        document.documentElement.setAttribute("data-theme", darkMode ? "dark" : "light")
        this.setState({ darkMode })
    }

    private handleLanguageChange() {
        this.setState({
            locale: this.buildLocale(),
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
        const { locale } = this.state

        return (
            <ErrorBoundary>
                <ConfigProvider theme={{ algorithm: ThemeContext.algorithm() }} locale={locale}>
                    <App>
                        <ThemedResources />
                        <AppContainer />
                    </App>
                </ConfigProvider>
            </ErrorBoundary>
        )
    }
}
