import React from "react"
import IntegrationService from "./IntegrationService"
import { IntegrationResponse } from "./model/IntegrationResponse"
import Preloader from "../../core/components/preloader/Preloader"
import AppShellContext from "../../core/components/shell/AppShellContext"
import IntegrationIcons from "./icon/IntegrationIcons"
import { Card, Flex, Image } from "antd"
import { SettingOutlined } from "@ant-design/icons"
import If from "../../core/components/flowcontrol/If"
import IntegrationConfigurationModal from "./component/IntegrationConfigurationModal"

interface IntegrationsPageState {
    loading: boolean
    integrations: IntegrationResponse[]
    integrationBeingChanged?: IntegrationResponse
}

export default class IntegrationsPage extends React.Component<any, IntegrationsPageState> {
    private readonly service: IntegrationService

    constructor(props: any) {
        super(props)
        this.service = new IntegrationService()
        this.state = {
            loading: true,
            integrationBeingChanged: undefined,
            integrations: [],
        }
    }

    private closeConfigurationModal(valuesChanged: boolean) {
        this.setState({
            integrationBeingChanged: undefined,
            loading: valuesChanged,
        })

        if (valuesChanged) this.fetchIntegrations()
    }

    private openConfigurationModal(integration: IntegrationResponse) {
        this.setState({ integrationBeingChanged: integration })
    }

    private renderIntegration(integration: IntegrationResponse) {
        const { id, enabled, name, description } = integration
        const icon = IntegrationIcons[id]
        return (
            <Card style={{ minWidth: 300 }}>
                <Flex style={{ flexGrow: 1 }}>
                    <Flex style={{ flexGrow: 1 }}>
                        <Card.Meta
                            avatar={<Image src={icon} preview={false} width={125} />}
                            title={name}
                            description={
                                <>
                                    <b>Status:</b> {enabled ? "Enabled" : "Disabled"}
                                    <p>{description}</p>
                                </>
                            }
                        />
                    </Flex>
                    <Flex justify="start" align="start">
                        <SettingOutlined
                            key="setting"
                            style={{ fontSize: 18 }}
                            onClick={() => this.openConfigurationModal(integration)}
                        />
                    </Flex>
                </Flex>
            </Card>
        )
    }

    private renderIntegrations() {
        const { integrations } = this.state
        return <>{integrations.map(integration => this.renderIntegration(integration))}</>
    }

    private fetchIntegrations() {
        this.service.getAll().then(integrations =>
            this.setState({
                loading: false,
                integrations,
            }),
        )
    }

    componentDidMount() {
        this.fetchIntegrations()
        AppShellContext.get().updateConfig({
            title: "Integrations",
            subtitle: "Configuration of the nginx ignition integration with third-party apps",
        })
    }

    render() {
        const { loading, integrationBeingChanged } = this.state
        if (loading) return <Preloader loading />

        return (
            <>
                {this.renderIntegrations()}
                <If condition={integrationBeingChanged !== undefined}>
                    {() => (
                        <IntegrationConfigurationModal
                            integrationId={integrationBeingChanged!!.id}
                            onClose={valuesChanged => this.closeConfigurationModal(valuesChanged)}
                        />
                    )}
                </If>
            </>
        )
    }
}
