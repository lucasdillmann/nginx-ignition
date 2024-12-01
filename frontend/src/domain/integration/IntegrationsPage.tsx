import React from "react"
import IntegrationService from "./IntegrationService"
import { IntegrationResponse } from "./model/IntegrationResponse"
import Preloader from "../../core/components/preloader/Preloader"
import AppShellContext from "../../core/components/shell/AppShellContext"

interface IntegrationsPageState {
    loading: boolean
    integrations: IntegrationResponse[]
}

export default class IntegrationsPage extends React.Component<any, IntegrationsPageState> {
    private readonly service: IntegrationService

    constructor(props: any) {
        super(props)
        this.service = new IntegrationService()
        this.state = {
            loading: true,
            integrations: [],
        }
    }

    componentDidMount() {
        this.service.getAll().then(integrations =>
            this.setState({
                loading: false,
                integrations,
            }),
        )

        AppShellContext.get().updateConfig({
            title: "Integrations",
            subtitle: "Configuration of the nginx ignition integration with third-party apps",
        })
    }

    render() {
        const { loading } = this.state
        if (loading) return <Preloader loading />

        return "TODO: Implement this"
    }
}
