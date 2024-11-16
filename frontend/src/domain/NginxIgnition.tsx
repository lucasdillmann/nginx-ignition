import React from 'react';
import ApiClientEventDispatcher from "../core/apiclient/event/ApiClientEventDispatcher";
import AuthenticationApiClientEventListener from "../core/authentication/AuthenticationApiClientEventListener";
import { App, ConfigProvider } from "antd";
import AppContext, { AppContextData, loadAppContextData } from "../core/components/context/AppContext";
import ErrorBoundary from "../core/components/errorboundary/ErrorBoundary";
import AppRouter from "../core/components/router/AppRouter";
import Routes from "./Routes";
import FullPagePreloader from "../core/components/preloader/FullPagePreloader";

interface NginxIgnitionState {
    context?: AppContextData
}

export default class NginxIgnition extends React.Component<unknown, NginxIgnitionState> {
    constructor(props: any) {
        super(props);
        this.state = {}
        ApiClientEventDispatcher.register(new AuthenticationApiClientEventListener())
    }

    componentDidMount() {
        loadAppContextData()
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
            return this.renderContainer(<FullPagePreloader />)

        return this.renderContainer(
            <AppContext.Provider value={context}>
                <AppRouter routes={Routes} />
            </AppContext.Provider>
        )
    }
}
