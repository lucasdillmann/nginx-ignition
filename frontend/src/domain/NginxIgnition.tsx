import React from 'react';
import { App, ConfigProvider } from "antd";
import ErrorBoundary from "../core/components/errorboundary/ErrorBoundary";
import AppContainer from "./AppContainer";
import ApiClientEventDispatcher from "../core/apiclient/event/ApiClientEventDispatcher";
import AuthenticationApiClientEventListener from "../core/authentication/AuthenticationApiClientEventListener";

export default class NginxIgnition extends React.PureComponent {
    componentDidMount() {
        ApiClientEventDispatcher.register(new AuthenticationApiClientEventListener())
        const preloader = document.getElementById('preloader') as HTMLElement
        preloader?.remove()
    }

    render() {
        return (
            <ErrorBoundary>
                <React.StrictMode>
                    <ConfigProvider>
                        <App>
                            <AppContainer />
                        </App>
                    </ConfigProvider>
                </React.StrictMode>
            </ErrorBoundary>
        )
    }
}
