import React from 'react';
import ApiClientEventDispatcher from "../core/apiclient/event/ApiClientEventDispatcher";
import AuthenticationApiClientEventListener from "../core/authentication/AuthenticationApiClientEventListener";
import { App, ConfigProvider } from "antd";
import ApplicationContext, { ApplicationContextData, startApplicationContext } from "../core/components/context/ApplicationContext";
import ErrorBoundary from "../core/components/errorboundary/ErrorBoundary";

interface NginxIgnitionState {
    context?: ApplicationContextData
}

export default class NginxIgnition extends React.Component<unknown, NginxIgnitionState> {
    constructor(props: any) {
        super(props);
        this.state = {}
        ApiClientEventDispatcher.register(new AuthenticationApiClientEventListener())
    }

    componentDidMount() {
        startApplicationContext()
            .then(context => {
                this.setState((current) => ({
                    ...current,
                    context,
                }))
            })

        const preloader = document.getElementById('preloader') as HTMLElement
        preloader.remove()
    }

    private renderContainer(root: React.ReactElement): React.ReactElement {
        return (
            <ErrorBoundary>
                <React.StrictMode>
                    <ConfigProvider>
                        <App>
                            {root}
                        </App>
                    </ConfigProvider>
                </React.StrictMode>
            </ErrorBoundary>
        );
    }

    render() {
        const { context } = this.state
        if (context == null)
            return this.renderContainer(
                <p>Loading</p>
            )

        return this.renderContainer(
            <ApplicationContext.Provider value={context}>
                <p>Hello there</p>
            </ApplicationContext.Provider>
        )
    }
}
