import React from 'react';
import AppContext, { AppContextData, loadAppContextData } from "../core/components/context/AppContext";
import AppRouter from "../core/components/router/AppRouter";
import Routes from "./Routes";
import FullPagePreloader from "../core/components/preloader/FullPagePreloader";
import FullPageError from "../core/components/error/FullPageError";

interface AppContainerState {
    context?: AppContextData
    error?: Error
}

export default class AppContainer extends React.Component<unknown, AppContainerState> {
    constructor(props: any) {
        super(props);
        this.state = {}
    }

    componentDidMount() {
        loadAppContextData()
            .then(context => this.setState({
                context,
                error: undefined,
            }))
            .catch(error => this.setState({
                error,
                context: undefined,
            }))
    }

    render() {
        const { context, error } = this.state
        if (error !== undefined)
            return <FullPageError error={error} />

        if (context === undefined)
            return <FullPagePreloader />

        return (
            <AppContext.Provider value={context}>
                <AppRouter routes={Routes} />
            </AppContext.Provider>
        )
    }
}
